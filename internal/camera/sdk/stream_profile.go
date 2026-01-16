//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/StreamProfile.h>
*/
import "C"
import "fmt"

// Close releases the stream profile list.
func (l *StreamProfileList) Close() error {
	if l == nil || l.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_stream_profile_list(l.ptr, &err)
	l.ptr = nil
	return CheckError(err)
}

// Count returns number of profiles in the list.
func (l *StreamProfileList) Count() (int, error) {
	if l == nil || l.ptr == nil {
		return 0, fmt.Errorf("profile list is nil")
	}
	var err *C.ob_error
	count := C.ob_stream_profile_list_count(l.ptr, &err)
	return int(count), CheckError(err)
}

// GetProfile returns the profile at index.
func (l *StreamProfileList) GetProfile(index int) (*StreamProfile, error) {
	if l == nil || l.ptr == nil {
		return nil, fmt.Errorf("profile list is nil")
	}
	var err *C.ob_error
	p := C.ob_stream_profile_list_get_profile(l.ptr, C.int(index), &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if p == nil {
		return nil, fmt.Errorf("profile at %d is nil", index)
	}
	return &StreamProfile{ptr: p}, nil
}

// Close releases the stream profile.
func (p *StreamProfile) Close() error {
	if p == nil || p.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_stream_profile(p.ptr, &err)
	p.ptr = nil
	return CheckError(err)
}

// Width returns video width.
func (p *StreamProfile) Width() (int, error) {
	if p == nil || p.ptr == nil {
		return 0, fmt.Errorf("stream profile is nil")
	}
	var err *C.ob_error
	w := C.ob_video_stream_profile_width(p.ptr, &err)
	return int(w), CheckError(err)
}

// Height returns video height.
func (p *StreamProfile) Height() (int, error) {
	if p == nil || p.ptr == nil {
		return 0, fmt.Errorf("stream profile is nil")
	}
	var err *C.ob_error
	h := C.ob_video_stream_profile_height(p.ptr, &err)
	return int(h), CheckError(err)
}

// FPS returns frames per second.
func (p *StreamProfile) FPS() (int, error) {
	if p == nil || p.ptr == nil {
		return 0, fmt.Errorf("stream profile is nil")
	}
	var err *C.ob_error
	fps := C.ob_video_stream_profile_fps(p.ptr, &err)
	return int(fps), CheckError(err)
}

// Format returns the frame format.
func (p *StreamProfile) Format() (FrameFormat, error) {
	if p == nil || p.ptr == nil {
		return 0, fmt.Errorf("stream profile is nil")
	}
	var err *C.ob_error
	f := C.ob_stream_profile_format(p.ptr, &err)
	return FrameFormat(f), CheckError(err)
}
