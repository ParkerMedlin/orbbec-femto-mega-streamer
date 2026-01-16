//go:build windows && cgo

package camera

import (
	"unsafe"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

func init() {
	pipelineFactory = func(dev *sdk.Device) (streamPipeline, error) {
		p, err := sdk.CreatePipelineWithDevice(dev)
		if err != nil {
			return nil, err
		}
		return &sdkPipelineAdapter{p: p}, nil
	}
}

// sdkPipelineAdapter adapts sdk.Pipeline to the streamPipeline interface.
type sdkPipelineAdapter struct {
	p *sdk.Pipeline
}

func (a *sdkPipelineAdapter) EnableStream(profile *sdk.StreamProfile) error {
	return a.p.EnableStream(profile)
}

func (a *sdkPipelineAdapter) Start() error {
	return a.p.Start()
}

func (a *sdkPipelineAdapter) Stop() error {
	return a.p.Stop()
}

func (a *sdkPipelineAdapter) WaitForFrameSet(timeoutMs int) (frameSet, error) {
	fs, err := a.p.WaitForFrameSet(timeoutMs)
	if err != nil || fs == nil {
		return nil, err
	}
	return &sdkFrameSetAdapter{fs: fs}, nil
}

func (a *sdkPipelineAdapter) Close() error {
	return a.p.Close()
}

type sdkFrameSetAdapter struct {
	fs *sdk.FrameSet
}

func (a *sdkFrameSetAdapter) DepthFrame() (frame, error) {
	f, err := a.fs.DepthFrame()
	if err != nil || f == nil {
		return nil, err
	}
	return &sdkFrameAdapter{f: f}, nil
}

func (a *sdkFrameSetAdapter) ColorFrame() (frame, error) {
	f, err := a.fs.ColorFrame()
	if err != nil || f == nil {
		return nil, err
	}
	return &sdkFrameAdapter{f: f}, nil
}

func (a *sdkFrameSetAdapter) IRFrame() (frame, error) {
	f, err := a.fs.IRFrame()
	if err != nil || f == nil {
		return nil, err
	}
	return &sdkFrameAdapter{f: f}, nil
}

func (a *sdkFrameSetAdapter) Close() error {
	return a.fs.Close()
}

type sdkFrameAdapter struct {
	f *sdk.Frame
}

func (a *sdkFrameAdapter) Close() error                     { return a.f.Close() }
func (a *sdkFrameAdapter) Index() (uint64, error)           { return a.f.Index() }
func (a *sdkFrameAdapter) Format() (sdk.FrameFormat, error) { return a.f.Format() }
func (a *sdkFrameAdapter) Timestamp() (uint64, error)       { return a.f.Timestamp() }
func (a *sdkFrameAdapter) SystemTimestamp() (uint64, error) { return a.f.SystemTimestamp() }
func (a *sdkFrameAdapter) Width() (int, error)              { return a.f.Width() }
func (a *sdkFrameAdapter) Height() (int, error)             { return a.f.Height() }
func (a *sdkFrameAdapter) DataSize() (int, error)           { return a.f.DataSize() }
func (a *sdkFrameAdapter) DataPtr() (unsafe.Pointer, error) { return a.f.DataPtr() }
