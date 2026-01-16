//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Device.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Count returns the number of detected devices.
func (d *DeviceList) Count() (int, error) {
	if d == nil || d.ptr == nil {
		return 0, nil
	}
	var err *C.ob_error
	count := C.ob_device_list_device_count(d.ptr, &err)
	if e := CheckError(err); e != nil {
		return 0, e
	}
	return int(count), nil
}

// SerialNumber returns the serial number for the device at index.
func (d *DeviceList) SerialNumber(index int) (string, error) {
	if d == nil || d.ptr == nil {
		return "", fmt.Errorf("device list is nil")
	}
	var err *C.ob_error
	serial := C.ob_device_list_get_device_serial_number(d.ptr, C.uint32_t(index), &err)
	if e := CheckError(err); e != nil {
		return "", e
	}
	return C.GoString(serial), nil
}

// ConnectionType returns the connection type string for the device at index.
func (d *DeviceList) ConnectionType(index int) (string, error) {
	if d == nil || d.ptr == nil {
		return "", fmt.Errorf("device list is nil")
	}
	var err *C.ob_error
	ct := C.ob_device_list_get_device_connection_type(d.ptr, C.uint32_t(index), &err)
	if e := CheckError(err); e != nil {
		return "", e
	}
	return C.GoString(ct), nil
}

// GetDevice returns a Device wrapper for the device at index.
func (d *DeviceList) GetDevice(index int) (*Device, error) {
	if d == nil || d.ptr == nil {
		return nil, fmt.Errorf("device list is nil")
	}
	var err *C.ob_error
	dev := C.ob_device_list_get_device(d.ptr, C.uint32_t(index), &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &Device{ptr: dev}, nil
}

// GetDeviceBySerial returns a Device by serial number.
func (d *DeviceList) GetDeviceBySerial(serial string) (*Device, error) {
	if d == nil || d.ptr == nil {
		return nil, fmt.Errorf("device list is nil")
	}
	cSerial := C.CString(serial)
	defer C.free(unsafe.Pointer(cSerial))
	var err *C.ob_error
	dev := C.ob_device_list_get_device_by_serial_number(d.ptr, cSerial, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &Device{ptr: dev}, nil
}

// Close releases the device list.
func (d *DeviceList) Close() error {
	if d == nil || d.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_device_list(d.ptr, &err)
	d.ptr = nil
	return CheckError(err)
}
