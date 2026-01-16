package camera

import (
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

// fake implementations to run tests without the real SDK.
type fakePipeline struct {
	frames     chan frameSet
	startCount int
	stopCount  int
	enabled    bool
	mu         sync.Mutex
	closed     bool
}

type testProvider struct{}

func (testProvider) Profiles(t StreamType) ([]StreamProfileInfo, error) {
	switch t {
	case StreamTypeDepth:
		return []StreamProfileInfo{{
			Width:  2,
			Height: 2,
			FPS:    30,
			Format: PixelFormatDepthU16,
		}}, nil
	case StreamTypeColor:
		return []StreamProfileInfo{{
			Width:  640,
			Height: 480,
			FPS:    15,
			Format: PixelFormatRGB8,
		}, {
			Width:  640,
			Height: 480,
			FPS:    30,
			Format: PixelFormatRGB8,
		}}, nil
	default:
		return nil, nil
	}
}

func newFakePipeline() *fakePipeline {
	return &fakePipeline{frames: make(chan frameSet, 4)}
}

func (p *fakePipeline) EnableStream(_ *sdk.StreamProfile) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = true
	return nil
}

func (p *fakePipeline) Start() error {
	p.mu.Lock()
	p.startCount++
	p.mu.Unlock()
	return nil
}

func (p *fakePipeline) Stop() error {
	p.mu.Lock()
	p.stopCount++
	p.closed = true
	p.mu.Unlock()
	return nil
}

func (p *fakePipeline) WaitForFrameSet(timeoutMs int) (frameSet, error) {
	select {
	case fs := <-p.frames:
		return fs, nil
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		return nil, nil
	}
}

func (p *fakePipeline) Close() error { return nil }

type fakeFrameSet struct {
	depth frame
	color frame
	ir    frame
}

func (f *fakeFrameSet) DepthFrame() (frame, error) { return f.depth, nil }
func (f *fakeFrameSet) ColorFrame() (frame, error) { return f.color, nil }
func (f *fakeFrameSet) IRFrame() (frame, error)    { return f.ir, nil }
func (f *fakeFrameSet) Close() error               { return nil }

type fakeFrame struct {
	data   []byte
	width  int
	height int
	format sdk.FrameFormat
	idx    uint64
	ts     uint64
}

func (f *fakeFrame) Close() error                     { return nil }
func (f *fakeFrame) Index() (uint64, error)           { return f.idx, nil }
func (f *fakeFrame) Format() (sdk.FrameFormat, error) { return f.format, nil }
func (f *fakeFrame) Timestamp() (uint64, error)       { return f.ts, nil }
func (f *fakeFrame) SystemTimestamp() (uint64, error) { return f.ts, nil }
func (f *fakeFrame) Width() (int, error)              { return f.width, nil }
func (f *fakeFrame) Height() (int, error)             { return f.height, nil }
func (f *fakeFrame) DataSize() (int, error)           { return len(f.data), nil }
func (f *fakeFrame) DataPtr() (unsafe.Pointer, error) {
	if len(f.data) == 0 {
		return nil, nil
	}
	return unsafe.Pointer(&f.data[0]), nil
}

// TestStreamStartAndStop verifies start/stop and frame delivery using fakes.
func TestStreamStartAndStop(t *testing.T) {
	oldFactory := pipelineFactory
	fp := newFakePipeline()
	pipelineFactory = func(dev *sdk.Device) (streamPipeline, error) { return fp, nil }
	defer func() { pipelineFactory = oldFactory }()

	dev := &Device{sdkDev: &sdk.Device{}, streams: make(map[StreamType]*Stream)}
	cfg := StreamConfig{Type: StreamTypeDepth, Enabled: true, Width: 2, Height: 2, FPS: 30, Format: PixelFormatDepthU16}
	stream := newStream(dev, cfg)
	stream.provider = testProvider{}
	stream.selectProfile = func(StreamConfig) (*sdk.StreamProfile, error) { return &sdk.StreamProfile{}, nil }

	if err := stream.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	// Send one frame.
	raw := &fakeFrame{data: []byte{1, 2, 3, 4}, width: 2, height: 1, format: sdk.FormatY16, idx: 7, ts: 1234}
	fp.frames <- &fakeFrameSet{depth: raw}

	select {
	case f := <-stream.Frames():
		if f.FrameID != 7 {
			t.Fatalf("unexpected frame id %d", f.FrameID)
		}
		if f.Size() != len(raw.data) {
			t.Fatalf("size mismatch got %d want %d", f.Size(), len(raw.data))
		}
		f.Release()
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for frame")
	}

	if err := stream.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if fp.startCount != 1 || fp.stopCount == 0 {
		t.Fatalf("expected start/stop to be called, got start %d stop %d", fp.startCount, fp.stopCount)
	}
}

// TestStreamSetConfig restarts when config changes while running.
func TestStreamSetConfig(t *testing.T) {
	oldFactory := pipelineFactory
	fp := newFakePipeline()
	pipelineFactory = func(dev *sdk.Device) (streamPipeline, error) { return fp, nil }
	defer func() { pipelineFactory = oldFactory }()

	dev := &Device{sdkDev: &sdk.Device{}, streams: make(map[StreamType]*Stream)}
	cfg := StreamConfig{Type: StreamTypeColor, Enabled: true, Width: 640, Height: 480, FPS: 15, Format: PixelFormatRGB8}
	stream := newStream(dev, cfg)
	stream.provider = testProvider{}
	stream.selectProfile = func(StreamConfig) (*sdk.StreamProfile, error) { return &sdk.StreamProfile{}, nil }

	if err := stream.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	newCfg := cfg
	newCfg.FPS = 30
	if err := stream.SetConfig(newCfg); err != nil {
		t.Fatalf("set config failed: %v", err)
	}

	if fp.startCount < 2 {
		t.Fatalf("expected restart, start count=%d", fp.startCount)
	}
}
