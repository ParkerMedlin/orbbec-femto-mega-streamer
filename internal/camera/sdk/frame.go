//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Frame.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Close releases the frameset (frames are separate objects).
func (fs *FrameSet) Close() error {
	if fs == nil || fs.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_frame(fs.ptr, &err)
	fs.ptr = nil
	return CheckError(err)
}

// DepthFrame returns the depth frame from the set.
func (fs *FrameSet) DepthFrame() (*Frame, error) {
	if fs == nil || fs.ptr == nil {
		return nil, fmt.Errorf("frameset is nil")
	}
	var err *C.ob_error
	f := C.ob_frameset_depth_frame(fs.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if f == nil {
		return nil, nil
	}
	return &Frame{ptr: f}, nil
}

// ColorFrame returns the color frame from the set.
func (fs *FrameSet) ColorFrame() (*Frame, error) {
	if fs == nil || fs.ptr == nil {
		return nil, fmt.Errorf("frameset is nil")
	}
	var err *C.ob_error
	f := C.ob_frameset_color_frame(fs.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if f == nil {
		return nil, nil
	}
	return &Frame{ptr: f}, nil
}

// IRFrame returns the infrared frame from the set.
func (fs *FrameSet) IRFrame() (*Frame, error) {
	if fs == nil || fs.ptr == nil {
		return nil, fmt.Errorf("frameset is nil")
	}
	var err *C.ob_error
	f := C.ob_frameset_ir_frame(fs.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if f == nil {
		return nil, nil
	}
	return &Frame{ptr: f}, nil
}

// Close releases the frame.
func (f *Frame) Close() error {
	if f == nil || f.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_frame(f.ptr, &err)
	f.ptr = nil
	return CheckError(err)
}

// Index returns the frame index.
func (f *Frame) Index() (uint64, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	idx := C.ob_frame_index(f.ptr, &err)
	return uint64(idx), CheckError(err)
}

// Format returns the frame format.
func (f *Frame) Format() (FrameFormat, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	format := C.ob_frame_format(f.ptr, &err)
	return FrameFormat(format), CheckError(err)
}

// Timestamp returns the device timestamp in microseconds.
func (f *Frame) Timestamp() (uint64, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	ts := C.ob_frame_time_stamp_us(f.ptr, &err)
	return uint64(ts), CheckError(err)
}

// SystemTimestamp returns the host timestamp in microseconds.
func (f *Frame) SystemTimestamp() (uint64, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	ts := C.ob_frame_system_time_stamp_us(f.ptr, &err)
	return uint64(ts), CheckError(err)
}

// Width returns the video frame width.
func (f *Frame) Width() (int, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	w := C.ob_video_frame_width(f.ptr, &err)
	return int(w), CheckError(err)
}

// Height returns the video frame height.
func (f *Frame) Height() (int, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	h := C.ob_video_frame_height(f.ptr, &err)
	return int(h), CheckError(err)
}

// DataSize returns the byte size of the frame payload.
func (f *Frame) DataSize() (int, error) {
	if f == nil || f.ptr == nil {
		return 0, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	size := C.ob_frame_data_size(f.ptr, &err)
	return int(size), CheckError(err)
}

// DataPtr returns a raw pointer to frame data.
// The pointer is valid until the frame is closed.
func (f *Frame) DataPtr() (unsafe.Pointer, error) {
	if f == nil || f.ptr == nil {
		return nil, fmt.Errorf("frame is nil")
	}
	var err *C.ob_error
	ptr := C.ob_frame_data(f.ptr, &err)
	return ptr, CheckError(err)
}
