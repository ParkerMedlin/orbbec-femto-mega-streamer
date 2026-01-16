//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/Error.h>
*/
import "C"
import "fmt"

// ExceptionType maps to ob_exception_type.
type ExceptionType int

const (
	ExceptionUnknown          ExceptionType = C.OB_EXCEPTION_TYPE_UNKNOWN
	ExceptionCameraDisconnect ExceptionType = C.OB_EXCEPTION_TYPE_CAMERA_DISCONNECTED
	ExceptionPlatform         ExceptionType = C.OB_EXCEPTION_TYPE_PLATFORM
	ExceptionInvalidValue     ExceptionType = C.OB_EXCEPTION_TYPE_INVALID_VALUE
	ExceptionWrongSequence    ExceptionType = C.OB_EXCEPTION_TYPE_WRONG_API_CALL_SEQUENCE
	ExceptionNotImplemented   ExceptionType = C.OB_EXCEPTION_TYPE_NOT_IMPLEMENTED
	ExceptionIO               ExceptionType = C.OB_EXCEPTION_TYPE_IO
	ExceptionMemory           ExceptionType = C.OB_EXCEPTION_TYPE_MEMORY
)

// SDKError represents an error returned by the Orbbec SDK.
type SDKError struct {
	Status        int
	Message       string
	ExceptionType ExceptionType
	Function      string
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("OrbbecSDK error in %s: %s (status=%d, type=%d)", e.Function, e.Message, e.Status, e.ExceptionType)
}

// CheckError converts an SDK error pointer into a Go error and frees the C memory.
func CheckError(err *C.ob_error) error {
	if err == nil {
		return nil
	}
	message := C.GoString(C.ob_error_message(err))
	status := int(C.ob_error_status(err))
	exType := ExceptionType(C.ob_error_exception_type(err))
	fn := C.GoString(C.ob_error_function(err))
	C.ob_delete_error(err)
	return &SDKError{
		Status:        status,
		Message:       message,
		ExceptionType: exType,
		Function:      fn,
	}
}
