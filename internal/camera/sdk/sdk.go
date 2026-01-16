//go:build windows && cgo

package sdk

/*
#cgo CFLAGS: -I${SRCDIR}/../../../orbbec-sdk/OrbbecSDK_v1.10.27/SDK/include
#cgo windows LDFLAGS: -L${SRCDIR}/../../../orbbec-sdk/OrbbecSDK_v1.10.27/SDK/lib -lOrbbecSDK

#include <libobsensor/ObSensor.h>
#include <stdlib.h>
*/
import "C"

// Context wraps ob_context.
type Context struct {
	ptr *C.ob_context
}

// DeviceList wraps ob_device_list.
type DeviceList struct {
	ptr *C.ob_device_list
}

// Device wraps ob_device.
type Device struct {
	ptr *C.ob_device
}

// Pipeline wraps ob_pipeline.
type Pipeline struct {
	ptr    *C.ob_pipeline
	config *C.ob_config
}

// FrameSet wraps ob_frameset.
type FrameSet struct {
	ptr *C.ob_frameset
}

// Frame wraps ob_frame.
type Frame struct {
	ptr *C.ob_frame
}

// StreamProfileList wraps ob_stream_profile_list.
type StreamProfileList struct {
	ptr *C.ob_stream_profile_list
}

// StreamProfile wraps ob_stream_profile.
type StreamProfile struct {
	ptr *C.ob_stream_profile
}
