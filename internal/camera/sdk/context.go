//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Context.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// LogSeverity maps to ob_log_severity.
type LogSeverity int

const (
	LogDebug LogSeverity = C.OB_LOG_SEVERITY_DEBUG
	LogInfo  LogSeverity = C.OB_LOG_SEVERITY_INFO
	LogWarn  LogSeverity = C.OB_LOG_SEVERITY_WARN
	LogError LogSeverity = C.OB_LOG_SEVERITY_ERROR
	LogFatal LogSeverity = C.OB_LOG_SEVERITY_FATAL
	LogOff   LogSeverity = C.OB_LOG_SEVERITY_OFF
)

// CreateContext initializes the SDK runtime.
func CreateContext() (*Context, error) {
	var err *C.ob_error
	ctx := C.ob_create_context(&err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &Context{ptr: ctx}, nil
}

// Close releases the SDK context.
func (c *Context) Close() error {
	if c == nil || c.ptr == nil {
		return nil
	}
	// Drop callback handle first to avoid dangling references.
	if c.callbackHandle != 0 {
		c.callbackHandle.Delete()
		c.callbackHandle = 0
	}
	var err *C.ob_error
	C.ob_delete_context(c.ptr, &err)
	c.ptr = nil
	return CheckError(err)
}

// SetLogLevel configures global SDK logging.
func (c *Context) SetLogLevel(level LogSeverity) error {
	if c == nil || c.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_set_logger_severity(C.ob_log_severity(level), &err)
	return CheckError(err)
}

// SetDeviceChangedCallback is implemented in Phase 3 (hot-plug).
// Keeping a stub here clarifies the API surface and avoids nil pointers when Phase 3 wires callbacks.
// func (c *Context) SetDeviceChangedCallback(...) error

// QueryDeviceList enumerates connected devices.
func (c *Context) QueryDeviceList() (*DeviceList, error) {
	if c == nil || c.ptr == nil {
		return nil, nil
	}
	var err *C.ob_error
	list := C.ob_query_device_list(c.ptr, &err)
	if e := CheckError(err); e != nil {
		return nil, e
	}
	return &DeviceList{ptr: list}, nil
}

// EnableNetDeviceEnumeration toggles network device discovery.
func (c *Context) EnableNetDeviceEnumeration(enable bool) error {
	if c == nil || c.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_enable_net_device_enumeration(c.ptr, C.bool(enable), &err)
	return CheckError(err)
}

// EnableDeviceClockSync enables host/device clock sync.
func (c *Context) EnableDeviceClockSync(intervalMs uint64) error {
	if c == nil || c.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_enable_device_clock_sync(c.ptr, C.uint64_t(intervalMs), &err)
	return CheckError(err)
}

// FreeIdleMemory instructs SDK to release idle frame buffers.
func (c *Context) FreeIdleMemory() error {
	if c == nil || c.ptr == nil {
		return nil
	}
	var err *C.ob_error
	C.ob_free_idle_memory(c.ptr, &err)
	return CheckError(err)
}

// LoadLicense loads a license from file.
func (c *Context) LoadLicense(path, key string) error {
	if c == nil || c.ptr == nil {
		return nil
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var err *C.ob_error
	C.ob_load_license(cPath, cKey, &err)
	return CheckError(err)
}
