//go:build !windows || !cgo

package sdk

import (
	"fmt"
	"unsafe"
)

// Minimal stub implementations so packages depending on sdk compile on
// non-Windows or non-CGo builds. All functions return an unsupported-platform
// error to make the limitation explicit at runtime.

type LogSeverity int

const (
	LogDebug LogSeverity = iota
	LogInfo
	LogWarn
	LogError
	LogFatal
	LogOff
)

type ExceptionType int

const (
	ExceptionUnknown ExceptionType = iota
	ExceptionCameraDisconnect
	ExceptionPlatform
	ExceptionInvalidValue
	ExceptionWrongSequence
	ExceptionNotImplemented
	ExceptionIO
	ExceptionMemory
)

type SensorType int

const (
	SensorUnknown SensorType = iota
	SensorIR
	SensorColor
	SensorDepth
	SensorAccel
	SensorGyro
)

type FrameFormat int

const (
	FormatYUYV FrameFormat = iota
	FormatNV12
	FormatMJPG
	FormatRGB
	FormatBGR
	FormatRGBA
	FormatBGRA
	FormatY16
	FormatY8
	FormatACCEL
	FormatGYRO
)

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

type Context struct{}

type DeviceList struct{}
type Device struct{}
type Pipeline struct{}
type FrameSet struct{}
type Frame struct{}
type StreamProfileList struct{}
type StreamProfile struct{}

// DeviceChangedCallback mirrors the Windows build signature.
type DeviceChangedCallback func(added []string, removed []string)

func unsupported() error {
	return fmt.Errorf("orbbec sdk is only available on Windows with cgo enabled")
}

func CreateContext() (*Context, error)                   { return nil, unsupported() }
func (c *Context) Close() error                          { return nil }
func (c *Context) SetLogLevel(LogSeverity) error         { return unsupported() }
func (c *Context) QueryDeviceList() (*DeviceList, error) { return nil, unsupported() }
func (c *Context) EnableNetDeviceEnumeration(bool) error { return unsupported() }
func (c *Context) EnableDeviceClockSync(uint64) error    { return unsupported() }
func (c *Context) FreeIdleMemory() error                 { return unsupported() }
func (c *Context) LoadLicense(string, string) error      { return unsupported() }
func (c *Context) SetDeviceChangedCallback(DeviceChangedCallback) (func(), error) {
	return func() {}, unsupported()
}

func (d *DeviceList) Count() (int, error)                       { return 0, unsupported() }
func (d *DeviceList) SerialNumber(int) (string, error)          { return "", unsupported() }
func (d *DeviceList) ConnectionType(int) (string, error)        { return "", unsupported() }
func (d *DeviceList) GetDevice(int) (*Device, error)            { return nil, unsupported() }
func (d *DeviceList) GetDeviceBySerial(string) (*Device, error) { return nil, unsupported() }
func (d *DeviceList) Close() error                              { return nil }

func (d *Device) Close() error                  { return unsupported() }
func (d *Device) GetInfo() (*DeviceInfo, error) { return nil, unsupported() }
func (d *Device) GetStreamProfileList(SensorType) (*StreamProfileList, error) {
	return nil, unsupported()
}

func (p *Pipeline) EnableStream(*StreamProfile) error      { return unsupported() }
func (p *Pipeline) Start() error                           { return unsupported() }
func (p *Pipeline) Stop() error                            { return unsupported() }
func (p *Pipeline) WaitForFrameSet(int) (*FrameSet, error) { return nil, unsupported() }
func (p *Pipeline) GetStreamProfileList(SensorType) (*StreamProfileList, error) {
	return nil, unsupported()
}
func (p *Pipeline) Close() error { return unsupported() }

func (f *FrameSet) Close() error                { return unsupported() }
func (f *FrameSet) DepthFrame() (*Frame, error) { return nil, unsupported() }
func (f *FrameSet) ColorFrame() (*Frame, error) { return nil, unsupported() }
func (f *FrameSet) IRFrame() (*Frame, error)    { return nil, unsupported() }

func (f *Frame) Close() error                     { return unsupported() }
func (f *Frame) Index() (uint64, error)           { return 0, unsupported() }
func (f *Frame) Format() (FrameFormat, error)     { return 0, unsupported() }
func (f *Frame) Timestamp() (int64, error)        { return 0, unsupported() }
func (f *Frame) SystemTimestamp() (int64, error)  { return 0, unsupported() }
func (f *Frame) Width() (int, error)              { return 0, unsupported() }
func (f *Frame) Height() (int, error)             { return 0, unsupported() }
func (f *Frame) DataSize() (int, error)           { return 0, unsupported() }
func (f *Frame) DataPtr() (unsafe.Pointer, error) { return nil, unsupported() }

func (l *StreamProfileList) Close() error        { return unsupported() }
func (l *StreamProfileList) Count() (int, error) { return 0, unsupported() }
func (l *StreamProfileList) GetProfile(int) (*StreamProfile, error) {
	return nil, unsupported()
}

func (p *StreamProfile) Close() error                 { return unsupported() }
func (p *StreamProfile) Width() (int, error)          { return 0, unsupported() }
func (p *StreamProfile) Height() (int, error)         { return 0, unsupported() }
func (p *StreamProfile) FPS() (int, error)            { return 0, unsupported() }
func (p *StreamProfile) Format() (FrameFormat, error) { return 0, unsupported() }
