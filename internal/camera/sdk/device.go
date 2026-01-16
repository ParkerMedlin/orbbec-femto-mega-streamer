//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Device.h>
#include <libobsensor/h/Sensor.h>
*/
import "C"
import "fmt"

// DeviceInfo contains basic metadata about the camera.
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

// Close releases the device handle.
func (d *Device) Close() error {
	if d == nil || d.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_delete_device(d.ptr, &err)
	d.ptr = nil
	return CheckError(err)
}

// GetInfo queries device metadata.
func (d *Device) GetInfo() (*DeviceInfo, error) {
	if d == nil || d.ptr == nil {
		return nil, fmt.Errorf("device is nil")
	}
	var err *C.ob_error
	info := C.ob_device_get_device_info(d.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if info == nil {
		return nil, fmt.Errorf("device info is nil")
	}
	defer func() {
		var derr *C.ob_error
		C.ob_delete_device_info(info, &derr)
		_ = CheckError(derr)
	}()

	getStr := func(fn func(*C.ob_device_info, **C.ob_error) *C.char) string {
		var innerErr *C.ob_error
		val := fn(info, &innerErr)
		_ = CheckError(innerErr)
		if val == nil {
			return ""
		}
		return C.GoString(val)
	}

	return &DeviceInfo{
		Name:            getStr(C.ob_device_info_name),
		Serial:          getStr(C.ob_device_info_serial_number),
		FirmwareVersion: getStr(C.ob_device_info_firmware_version),
		ConnectionType:  getStr(C.ob_device_info_connection_type),
		IPAddress:       getStr(C.ob_device_info_ip_address),
		HardwareVersion: getStr(C.ob_device_info_hardware_version),
		USBType:         getStr(C.ob_device_info_usb_type),
		ASICName:        getStr(C.ob_device_info_asicName),
	}, nil
}

// GetStreamProfileList fetches supported profiles for a given sensor type.
func (d *Device) GetStreamProfileList(sensor SensorType) (*StreamProfileList, error) {
	if d == nil || d.ptr == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Acquire sensor list.
	var err *C.ob_error
	sensors := C.ob_device_get_sensor_list(d.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	if sensors == nil {
		return nil, fmt.Errorf("sensor list is nil")
	}
	defer func() {
		var serr *C.ob_error
		C.ob_delete_sensor_list(sensors, &serr)
		_ = CheckError(serr)
	}()

	// Grab sensor by type.
	var sErr *C.ob_error
	s := C.ob_sensor_list_get_sensor_by_type(sensors, C.ob_sensor_type(sensor), &sErr)
	if e := CheckError(sErr); e != nil {
		return nil, e
	}
	if s == nil {
		return nil, fmt.Errorf("sensor type %d not found", sensor)
	}
	defer func() {
		var derr *C.ob_error
		C.ob_delete_sensor(s, &derr)
		_ = CheckError(derr)
	}()

	// Fetch profile list.
	var pErr *C.ob_error
	list := C.ob_sensor_get_stream_profile_list(s, &pErr)
	if e := CheckError(pErr); e != nil {
		return nil, e
	}
	return &StreamProfileList{ptr: list}, nil
}
