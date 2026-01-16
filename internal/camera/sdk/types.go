//go:build windows && cgo

package sdk

/*
#include <libobsensor/h/ObTypes.h>
*/
import "C"

// SensorType maps to ob_sensor_type.
type SensorType int

const (
	SensorUnknown SensorType = C.OB_SENSOR_UNKNOWN
	SensorIR      SensorType = C.OB_SENSOR_IR
	SensorColor   SensorType = C.OB_SENSOR_COLOR
	SensorDepth   SensorType = C.OB_SENSOR_DEPTH
	SensorAccel   SensorType = C.OB_SENSOR_ACCEL
	SensorGyro    SensorType = C.OB_SENSOR_GYRO
)

// FrameFormat maps to ob_format.
type FrameFormat int

const (
	FormatYUYV  FrameFormat = C.OB_FORMAT_YUYV
	FormatNV12  FrameFormat = C.OB_FORMAT_NV12
	FormatMJPG  FrameFormat = C.OB_FORMAT_MJPG
	FormatRGB   FrameFormat = C.OB_FORMAT_RGB
	FormatBGR   FrameFormat = C.OB_FORMAT_BGR
	FormatRGBA  FrameFormat = C.OB_FORMAT_RGBA
	FormatBGRA  FrameFormat = C.OB_FORMAT_BGRA
	FormatY16   FrameFormat = C.OB_FORMAT_Y16
	FormatY8    FrameFormat = C.OB_FORMAT_Y8
	FormatACCEL FrameFormat = C.OB_FORMAT_ACCEL
	FormatGYRO  FrameFormat = C.OB_FORMAT_GYRO
)
