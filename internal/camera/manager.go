package camera

import (
	"context"
	"fmt"
	"sync"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

// DeviceInfo mirrors SDK metadata but keeps the camera package decoupled.
type DeviceInfo struct {
	Name            string
	Serial          string
	FirmwareVersion string
	ConnectionType  string
	IPAddress       string
	HardwareVersion string
	USBType         string
	ASICName        string
}

// Device wraps an SDK device handle.
type Device struct {
	serial  string
	sdkDev  *sdk.Device
	info    DeviceInfo
	streams map[StreamType]*Stream

	mu     sync.Mutex
	closed bool
}

// Info returns cached device metadata.
func (d *Device) Info() DeviceInfo {
	if d == nil {
		return DeviceInfo{}
	}
	return d.info
}

// Close releases the underlying SDK device.
func (d *Device) Close() error {
	if d == nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil
	}
	d.closed = true

	for _, s := range d.streams {
		_ = s.Stop()
	}
	d.streams = nil

	if d.sdkDev != nil {
		return d.sdkDev.Close()
	}
	return nil
}

// Profiles implements ProfileProvider by querying supported profiles from the device.
func (d *Device) Profiles(streamType StreamType) ([]StreamProfileInfo, error) {
	return d.profiles(streamType)
}

// DeviceManager coordinates discovery, opening, and hot-plug events.
type DeviceManager struct {
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	sdk        sdkAdapter
	devices    map[string]*Device
	added      chan string
	removed    chan string
	unregister func()
	closed     bool
}

// sdkAdapter abstracts the SDK for testability.
type sdkAdapter interface {
	ListDevices() ([]DeviceInfo, error)
	OpenDevice(serial string) (*sdk.Device, *DeviceInfo, error)
	SetDeviceChangedCallback(cb sdk.DeviceChangedCallback) (func(), error)
	Close() error
}

// defaultSDKAdapter is the production implementation.
type defaultSDKAdapter struct {
	ctx *sdk.Context
}

// newSDKAdapter can be swapped in tests.
var newSDKAdapter = func() (sdkAdapter, error) {
	ctx, err := sdk.CreateContext()
	if err != nil {
		return nil, err
	}
	return &defaultSDKAdapter{ctx: ctx}, nil
}

// NewDeviceManager initializes the SDK context and registers hot-plug callbacks.
func NewDeviceManager(ctx context.Context) (*DeviceManager, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cctx, cancel := context.WithCancel(ctx)

	sdkImpl, err := newSDKAdapter()
	if err != nil {
		cancel()
		return nil, err
	}

	mgr := &DeviceManager{
		ctx:     cctx,
		cancel:  cancel,
		sdk:     sdkImpl,
		devices: make(map[string]*Device),
		added:   make(chan string, 8),
		removed: make(chan string, 8),
	}

	unregister, err := sdkImpl.SetDeviceChangedCallback(func(added []string, removed []string) {
		for _, s := range added {
			mgr.handleDeviceAdded(s)
		}
		for _, s := range removed {
			mgr.handleDeviceRemoved(s)
		}
	})
	if err != nil {
		cancel()
		_ = sdkImpl.Close()
		return nil, err
	}
	mgr.unregister = unregister

	return mgr, nil
}

// Discover enumerates connected devices and returns their info.
func (m *DeviceManager) Discover() ([]DeviceInfo, error) {
	if m == nil {
		return nil, fmt.Errorf("manager is nil")
	}
	return m.sdk.ListDevices()
}

// Open returns an opened device by serial. If serial is empty, the first available device is opened.
func (m *DeviceManager) Open(serial string) (*Device, error) {
	if m == nil {
		return nil, fmt.Errorf("manager is nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, fmt.Errorf("manager is closed")
	}

	if serial == "" {
		devices, err := m.sdk.ListDevices()
		if err != nil {
			return nil, err
		}
		if len(devices) == 0 {
			return nil, fmt.Errorf("no devices available")
		}
		serial = devices[0].Serial
	}

	if existing, ok := m.devices[serial]; ok {
		return existing, nil
	}

	sdkDev, info, err := m.sdk.OpenDevice(serial)
	if err != nil {
		return nil, &DeviceError{Serial: serial, Op: "open", Err: err}
	}
	dev := &Device{
		serial:  serial,
		sdkDev:  sdkDev,
		info:    *info,
		streams: make(map[StreamType]*Stream),
	}
	m.devices[serial] = dev
	return dev, nil
}

// Close releases all devices and the SDK context.
func (m *DeviceManager) Close() error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.mu.Unlock()

	if m.unregister != nil {
		m.unregister()
	}
	if m.cancel != nil {
		m.cancel()
	}

	m.mu.Lock()
	for serial, dev := range m.devices {
		_ = dev.Close()
		delete(m.devices, serial)
	}
	m.mu.Unlock()

	if m.sdk != nil {
		_ = m.sdk.Close()
	}

	close(m.added)
	close(m.removed)
	return nil
}

// OnDeviceAdded returns a channel of serial numbers for hot-plug events.
func (m *DeviceManager) OnDeviceAdded() <-chan string {
	return m.added
}

// OnDeviceRemoved returns a channel of serial numbers for unplug events.
func (m *DeviceManager) OnDeviceRemoved() <-chan string {
	return m.removed
}

func (m *DeviceManager) handleDeviceAdded(serial string) {
	m.mu.RLock()
	closed := m.closed
	m.mu.RUnlock()
	if closed {
		return
	}
	select {
	case m.added <- serial:
	default:
	}
}

func (m *DeviceManager) handleDeviceRemoved(serial string) {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	if dev, ok := m.devices[serial]; ok {
		_ = dev.Close()
		delete(m.devices, serial)
	}
	m.mu.Unlock()

	select {
	case m.removed <- serial:
	default:
	}
}

// ListDevices enumerates devices through the SDK adapter.
func (a *defaultSDKAdapter) ListDevices() ([]DeviceInfo, error) {
	list, err := a.ctx.QueryDeviceList()
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer list.Close()

	count, err := list.Count()
	if err != nil {
		return nil, err
	}

	devices := make([]DeviceInfo, 0, count)
	for i := 0; i < count; i++ {
		dev, err := list.GetDevice(i)
		if err != nil {
			return nil, err
		}
		info, err := dev.GetInfo()
		_ = dev.Close()
		if err != nil {
			return nil, err
		}
		devices = append(devices, sdkInfoToCamera(info))
	}
	return devices, nil
}

// OpenDevice opens a specific device by serial.
func (a *defaultSDKAdapter) OpenDevice(serial string) (*sdk.Device, *DeviceInfo, error) {
	list, err := a.ctx.QueryDeviceList()
	if err != nil {
		return nil, nil, err
	}
	if list == nil {
		return nil, nil, fmt.Errorf("no devices found")
	}
	defer list.Close()

	var dev *sdk.Device
	if serial == "" {
		dev, err = list.GetDevice(0)
	} else {
		dev, err = list.GetDeviceBySerial(serial)
	}
	if err != nil {
		return nil, nil, err
	}

	info, err := dev.GetInfo()
	if err != nil {
		_ = dev.Close()
		return nil, nil, err
	}
	cameraInfo := sdkInfoToCamera(info)
	return dev, &cameraInfo, nil
}

func (a *defaultSDKAdapter) SetDeviceChangedCallback(cb sdk.DeviceChangedCallback) (func(), error) {
	return a.ctx.SetDeviceChangedCallback(cb)
}

func (a *defaultSDKAdapter) Close() error {
	if a.ctx == nil {
		return nil
	}
	return a.ctx.Close()
}

func sdkInfoToCamera(info *sdk.DeviceInfo) DeviceInfo {
	if info == nil {
		return DeviceInfo{}
	}
	return DeviceInfo{
		Name:            info.Name,
		Serial:          info.Serial,
		FirmwareVersion: info.FirmwareVersion,
		ConnectionType:  info.ConnectionType,
		IPAddress:       info.IPAddress,
		HardwareVersion: info.HardwareVersion,
		USBType:         info.USBType,
		ASICName:        info.ASICName,
	}
}

// Stream returns (and lazily creates) a stream for the given type.
func (d *Device) Stream(t StreamType) *Stream {
	if d == nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil
	}
	if s, ok := d.streams[t]; ok {
		return s
	}

	var cfg StreamConfig
	switch t {
	case StreamTypeDepth:
		cfg = DefaultDepthConfig()
	case StreamTypeColor:
		cfg = DefaultColorConfig()
	case StreamTypeIR:
		cfg = StreamConfig{Type: StreamTypeIR, Enabled: true, Width: 640, Height: 480, FPS: 30, Format: PixelFormatIR16}
	default:
		cfg = StreamConfig{Type: t, Enabled: true}
	}

	s := newStream(d, cfg)
	d.streams[t] = s
	return s
}

// profiles queries supported profiles for a stream type.
func (d *Device) profiles(t StreamType) ([]StreamProfileInfo, error) {
	if d == nil || d.sdkDev == nil {
		return nil, fmt.Errorf("device is nil")
	}
	sensor, ok := sensorTypeFromStream(t)
	if !ok {
		return nil, fmt.Errorf("unsupported stream type %v", t)
	}

	list, err := d.sdkDev.GetStreamProfileList(sensor)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, fmt.Errorf("profile list is nil")
	}
	defer list.Close()

	count, err := list.Count()
	if err != nil {
		return nil, err
	}

	profiles := make([]StreamProfileInfo, 0, count)
	for i := 0; i < count; i++ {
		p, err := list.GetProfile(i)
		if err != nil {
			return nil, err
		}
		w, _ := p.Width()
		h, _ := p.Height()
		fps, _ := p.FPS()
		fmtRaw, _ := p.Format()
		_ = p.Close()

		pf, ok := pixelFormatFromSDK(fmtRaw)
		if !ok {
			continue
		}
		profiles = append(profiles, StreamProfileInfo{
			Width:  uint16(w),
			Height: uint16(h),
			FPS:    uint8(fps),
			Format: pf,
		})
	}
	return profiles, nil
}

// findProfile returns the SDK profile matching the provided config.
func (d *Device) findProfile(cfg StreamConfig) (*sdk.StreamProfile, error) {
	if d == nil || d.sdkDev == nil {
		return nil, fmt.Errorf("device is nil")
	}
	sensor, ok := sensorTypeFromStream(cfg.Type)
	if !ok {
		return nil, &ConfigError{Stream: cfg.Type, Field: "type", Value: cfg.Type, Reason: "unsupported"}
	}
	targetFormat, ok := pixelFormatToSDK(cfg.Format)
	if !ok {
		return nil, &ConfigError{Stream: cfg.Type, Field: "format", Value: cfg.Format, Reason: "unsupported format"}
	}

	list, err := d.sdkDev.GetStreamProfileList(sensor)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, fmt.Errorf("profile list is nil")
	}
	defer list.Close()

	count, err := list.Count()
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		p, err := list.GetProfile(i)
		if err != nil {
			return nil, err
		}
		w, _ := p.Width()
		h, _ := p.Height()
		fps, _ := p.FPS()
		fmtRaw, _ := p.Format()
		if w == int(cfg.Width) && h == int(cfg.Height) && fps == int(cfg.FPS) && fmtRaw == targetFormat {
			return p, nil // caller takes ownership
		}
		_ = p.Close()
	}
	return nil, &ConfigError{Stream: cfg.Type, Field: "profile", Value: cfg, Reason: "matching profile not found"}
}
