package camera

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

// frame abstracts over SDK frame implementations for testability.
type frame interface {
	Close() error
	Index() (uint64, error)
	Format() (sdk.FrameFormat, error)
	Timestamp() (uint64, error)
	SystemTimestamp() (uint64, error)
	Width() (int, error)
	Height() (int, error)
	DataSize() (int, error)
	DataPtr() (unsafe.Pointer, error)
}

// frameSet abstracts an SDK frameset for testability.
type frameSet interface {
	DepthFrame() (frame, error)
	ColorFrame() (frame, error)
	IRFrame() (frame, error)
	Close() error
}

// streamPipeline is the minimal pipeline API the Stream needs.
type streamPipeline interface {
	EnableStream(profile *sdk.StreamProfile) error
	Start() error
	Stop() error
	WaitForFrameSet(timeoutMs int) (frameSet, error)
	Close() error
}

// pipelineFactory creates a pipeline for a given device (set per-platform).
var pipelineFactory func(dev *sdk.Device) (streamPipeline, error)

// Stream wraps an SDK pipeline and delivers frames on a channel.
type Stream struct {
	device *Device

	mu            sync.RWMutex
	cfg           StreamConfig
	frames        chan *Frame
	pool          *framePool
	running       atomic.Bool
	cancel        context.CancelFunc
	done          chan struct{}
	pipe          streamPipeline
	captured      atomic.Uint64
	dropped       atomic.Uint64
	totalLatency  atomic.Int64
	lastFrameID   atomic.Uint64
	lastTimestamp atomic.Int64

	// test hooks
	provider      ProfileProvider
	selectProfile func(StreamConfig) (*sdk.StreamProfile, error)
}

func newStream(device *Device, cfg StreamConfig) *Stream {
	return &Stream{
		device: device,
		cfg:    cfg,
		frames: make(chan *Frame, 8),
		pool:   newFramePool(4, estimateBufferSize(cfg)),
		done:   make(chan struct{}),
	}
}

// Frames returns the channel that carries captured frames.
func (s *Stream) Frames() <-chan *Frame {
	if s == nil {
		return nil
	}
	return s.frames
}

// Config returns the current configuration.
func (s *Stream) Config() StreamConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

// Start begins frame capture if not already running.
func (s *Stream) Start() error {
	if s == nil {
		return fmt.Errorf("stream is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running.Load() {
		return nil
	}
	if s.device == nil || s.device.sdkDev == nil {
		return fmt.Errorf("device is nil")
	}

	// Validate config.
	provider := s.provider
	if provider == nil {
		provider = s.device
	}
	if err := s.cfg.Validate(provider); err != nil {
		return err
	}

	// Select profile.
	var selector func(StreamConfig) (*sdk.StreamProfile, error)
	if s.selectProfile != nil {
		selector = s.selectProfile
	} else {
		selector = func(cfg StreamConfig) (*sdk.StreamProfile, error) {
			return s.device.findProfile(cfg)
		}
	}
	profile, err := selector(s.cfg)
	if err != nil {
		return err
	}
	defer profile.Close()

	pipe, err := pipelineFactory(s.device.sdkDev)
	if err != nil {
		return err
	}
	if pipe == nil {
		return fmt.Errorf("pipeline factory returned nil")
	}
	if err := pipe.EnableStream(profile); err != nil {
		_ = pipe.Close()
		return err
	}
	if err := pipe.Start(); err != nil {
		_ = pipe.Close()
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.done = make(chan struct{})
	s.pipe = pipe
	s.running.Store(true)

	go s.captureLoop(ctx)
	return nil
}

// Stop halts capture and releases resources.
func (s *Stream) Stop() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	if !s.running.Load() {
		s.mu.Unlock()
		return nil
	}
	cancel := s.cancel
	done := s.done
	pipe := s.pipe
	s.cancel = nil
	s.pipe = nil
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}
	if pipe != nil {
		_ = pipe.Stop()
		_ = pipe.Close()
	}
	return nil
}

// SetConfig updates the stream configuration, restarting if already running.
func (s *Stream) SetConfig(cfg StreamConfig) error {
	if s == nil {
		return fmt.Errorf("stream is nil")
	}

	s.mu.Lock()
	s.cfg = cfg
	wasRunning := s.running.Load()
	s.mu.Unlock()

	if !wasRunning {
		return nil
	}
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

// IsRunning reports whether the stream is actively capturing.
func (s *Stream) IsRunning() bool {
	if s == nil {
		return false
	}
	return s.running.Load()
}

func (s *Stream) captureLoop(ctx context.Context) {
	defer func() {
		s.running.Store(false)
		close(s.done)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		start := time.Now()
		fs, err := s.pipe.WaitForFrameSet(1000)
		if err != nil || fs == nil {
			continue
		}

		var raw frame
		switch s.cfg.Type {
		case StreamTypeDepth:
			raw, err = fs.DepthFrame()
		case StreamTypeColor:
			raw, err = fs.ColorFrame()
		case StreamTypeIR:
			raw, err = fs.IRFrame()
		default:
			_ = fs.Close()
			continue
		}

		if err != nil || raw == nil {
			_ = fs.Close()
			continue
		}

		f, err := convertFrame(raw, s.cfg.Type, s.pool)
		_ = raw.Close()
		_ = fs.Close()
		if err != nil {
			continue
		}

		latency := time.Since(start)

		select {
		case s.frames <- f:
			s.recordCapture(f, latency)
		default:
			// Drop if channel full to avoid blocking capture loop.
			f.Release()
			s.recordDrop()
		}
	}
}

func convertFrame(src frame, streamType StreamType, pool *framePool) (*Frame, error) {
	if src == nil {
		return nil, fmt.Errorf("frame is nil")
	}
	width, err := src.Width()
	if err != nil {
		return nil, err
	}
	height, err := src.Height()
	if err != nil {
		return nil, err
	}
	size, err := src.DataSize()
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		return nil, fmt.Errorf("frame has no data")
	}
	formatRaw, err := src.Format()
	if err != nil {
		return nil, err
	}
	format, ok := pixelFormatFromSDK(formatRaw)
	if !ok {
		return nil, fmt.Errorf("unsupported frame format %v", formatRaw)
	}
	frameID, err := src.Index()
	if err != nil {
		return nil, err
	}
	ts, err := src.SystemTimestamp()
	if err != nil {
		return nil, err
	}
	ptr, err := src.DataPtr()
	if err != nil {
		return nil, err
	}

	f := pool.Get()
	if cap(f.Data) < size {
		f.Data = make([]byte, size)
	} else {
		f.Data = f.Data[:size]
	}
	// Copy from C buffer into Go slice.
	copy(f.Data, unsafe.Slice((*byte)(ptr), size))

	f.StreamType = streamType
	f.FrameID = frameID
	f.Timestamp = int64(ts * uint64(time.Microsecond/time.Nanosecond))
	f.Width = uint16(width)
	f.Height = uint16(height)
	f.Format = format
	return f, nil
}

func estimateBufferSize(cfg StreamConfig) int {
	pixelBytes := cfg.Format.BytesPerPixel()
	if pixelBytes == 0 {
		// Assume 2 bytes per pixel as a safe default for compressed/unknown formats.
		pixelBytes = 2
	}
	if cfg.Width == 0 || cfg.Height == 0 {
		return 1024 * 1024 // 1MB default
	}
	return int(cfg.Width) * int(cfg.Height) * pixelBytes
}

func (s *Stream) recordCapture(f *Frame, latency time.Duration) {
	if s == nil || f == nil {
		return
	}
	s.captured.Add(1)
	s.totalLatency.Add(latency.Nanoseconds())
	s.lastFrameID.Store(f.FrameID)
	s.lastTimestamp.Store(f.Timestamp)
}

func (s *Stream) recordDrop() {
	if s == nil {
		return
	}
	s.dropped.Add(1)
}
