//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Context.h>
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// DeviceChangedCallback receives serial numbers for added and removed devices.
type DeviceChangedCallback func(added []string, removed []string)

//export goDeviceChangedCallback
func goDeviceChangedCallback(removed *C.ob_device_list, added *C.ob_device_list, userData unsafe.Pointer) {
	if userData == nil {
		return
	}
	handle := cgo.Handle(uintptr(userData))
	cb, ok := handle.Value().(DeviceChangedCallback)
	if !ok || cb == nil {
		return
	}
	addedSerials := convertDeviceListToSerials(added)
	removedSerials := convertDeviceListToSerials(removed)
	cb(addedSerials, removedSerials)
}

func convertDeviceListToSerials(list *C.ob_device_list) []string {
	if list == nil {
		return nil
	}
	var err *C.ob_error
	count := int(C.ob_device_list_device_count(list, &err))
	if e := CheckError(err); e != nil {
		return nil
	}
	serials := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var sErr *C.ob_error
		serial := C.ob_device_list_get_device_serial_number(list, C.uint32_t(i), &sErr)
		if e := CheckError(sErr); e != nil {
			continue
		}
		if serial != nil {
			serials = append(serials, C.GoString(serial))
		}
	}
	return serials
}

// SetDeviceChangedCallback registers the hot-plug callback.
// It returns a cleanup function that removes the callback and frees internal handles.
func (c *Context) SetDeviceChangedCallback(cb DeviceChangedCallback) (func(), error) {
	if c == nil || c.ptr == nil {
		return func() {}, nil
	}

	// Clear any previous handle.
	if c.callbackHandle != 0 {
		c.callbackHandle.Delete()
		c.callbackHandle = 0
	}

	if cb == nil {
		var err *C.ob_error
		C.ob_set_device_changed_callback(c.ptr, nil, nil, &err)
		return func() {}, CheckError(err)
	}

	handle := cgo.NewHandle(cb)
	c.callbackHandle = handle

	var err *C.ob_error
	C.ob_set_device_changed_callback(
		c.ptr,
		(C.ob_device_changed_callback)(C.goDeviceChangedCallback),
		unsafe.Pointer(handle),
		&err,
	)
	if e := CheckError(err); e != nil {
		handle.Delete()
		c.callbackHandle = 0
		return func() {}, e
	}

	cleanup := func() {
		if c.callbackHandle != 0 {
			c.callbackHandle.Delete()
			c.callbackHandle = 0
		}
		var cerr *C.ob_error
		C.ob_set_device_changed_callback(c.ptr, nil, nil, &cerr)
		_ = CheckError(cerr)
	}
	return cleanup, nil
}
