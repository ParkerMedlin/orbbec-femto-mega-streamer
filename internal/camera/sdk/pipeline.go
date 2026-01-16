//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Pipeline.h>
*/
import "C"
import "fmt"

// CreatePipelineWithDevice creates a pipeline bound to the provided device.
func CreatePipelineWithDevice(dev *Device) (*Pipeline, error) {
	if dev == nil || dev.ptr == nil {
		return nil, fmt.Errorf("device is nil")
	}
	var err *C.ob_error
	p := C.ob_create_pipeline_with_device(dev.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &Pipeline{ptr: p}, nil
}

// EnableStream attaches a stream profile to the pipeline config.
// Config is created lazily on first call.
func (p *Pipeline) EnableStream(profile *StreamProfile) error {
	if p == nil || p.ptr == nil {
		return fmt.Errorf("pipeline is nil")
	}
	if profile == nil || profile.ptr == nil {
		return fmt.Errorf("stream profile is nil")
	}

	if p.config == nil {
		var cerr *C.ob_error
		p.config = C.ob_create_config(&cerr)
		if e := CheckError(cerr); e != nil {
			return e
		}
	}

	var err *C.ob_error
	C.ob_config_enable_stream(p.config, profile.ptr, &err)
	return CheckError(err)
}

// Start begins streaming. If no config has been built, starts with defaults.
func (p *Pipeline) Start() error {
	if p == nil || p.ptr == nil {
		return fmt.Errorf("pipeline is nil")
	}
	var err *C.ob_error
	if p.config != nil {
		C.ob_pipeline_start_with_config(p.ptr, p.config, &err)
	} else {
		C.ob_pipeline_start(p.ptr, &err)
	}
	return CheckError(err)
}

// Stop stops streaming.
func (p *Pipeline) Stop() error {
	if p == nil || p.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_pipeline_stop(p.ptr, &err)
	return CheckError(err)
}

// WaitForFrameSet blocks until a frameset arrives or timeout (ms). Returns nil on timeout.
func (p *Pipeline) WaitForFrameSet(timeoutMs int) (*FrameSet, error) {
	if p == nil || p.ptr == nil {
		return nil, fmt.Errorf("pipeline is nil")
	}
	var err *C.ob_error
	fs := C.ob_pipeline_wait_for_frameset(p.ptr, C.uint32_t(timeoutMs), &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if fs == nil {
		return nil, nil
	}
	return &FrameSet{ptr: fs}, nil
}

// GetStreamProfileList returns supported profiles for the given sensor on this pipeline.
func (p *Pipeline) GetStreamProfileList(sensor SensorType) (*StreamProfileList, error) {
	if p == nil || p.ptr == nil {
		return nil, fmt.Errorf("pipeline is nil")
	}
	var err *C.ob_error
	list := C.ob_pipeline_get_stream_profile_list(p.ptr, C.ob_sensor_type(sensor), &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &StreamProfileList{ptr: list}, nil
}

// Close releases pipeline and config.
func (p *Pipeline) Close() error {
	if p == nil {
		return nil
	}
	var firstErr error

	if p.config != nil {
		var err *C.ob_error
		C.ob_delete_config(p.config, &err)
		if e := CheckError(err); e != nil && firstErr == nil {
			firstErr = e
		}
		p.config = nil
	}

	// Best-effort stop before destroying pipeline.
	_ = p.Stop()

	if p.ptr != nil {
		var err *C.ob_error
		C.ob_delete_pipeline(p.ptr, &err)
		if e := CheckError(err); e != nil && firstErr == nil {
			firstErr = e
		}
		p.ptr = nil
	}

	return firstErr
}
