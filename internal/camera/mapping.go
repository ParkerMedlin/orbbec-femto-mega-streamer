package camera

import "github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"

// sensorTypeFromStream converts a StreamType to the corresponding SDK sensor type.
func sensorTypeFromStream(t StreamType) (sdk.SensorType, bool) {
	switch t {
	case StreamTypeDepth:
		return sdk.SensorDepth, true
	case StreamTypeColor:
		return sdk.SensorColor, true
	case StreamTypeIR:
		return sdk.SensorIR, true
	case StreamTypeIMU:
		// IMU is exposed as accelerometer/gyro; we map to accel by default.
		return sdk.SensorAccel, true
	default:
		return sdk.SensorUnknown, false
	}
}

// pixelFormatToSDK converts public PixelFormat to SDK FrameFormat.
func pixelFormatToSDK(p PixelFormat) (sdk.FrameFormat, bool) {
	switch p {
	case PixelFormatDepthU16, PixelFormatIR16:
		return sdk.FormatY16, true
	case PixelFormatRGB8:
		return sdk.FormatRGB, true
	case PixelFormatRGBA8:
		return sdk.FormatRGBA, true
	case PixelFormatMJPEG:
		return sdk.FormatMJPG, true
	case PixelFormatYUV422:
		return sdk.FormatYUYV, true
	case PixelFormatDepthF32:
		// No direct SDK format; not currently supported.
		return 0, false
	case PixelFormatIMU:
		return sdk.FormatACCEL, true
	default:
		return 0, false
	}
}

// pixelFormatFromSDK converts SDK FrameFormat to public PixelFormat.
func pixelFormatFromSDK(f sdk.FrameFormat) (PixelFormat, bool) {
	switch f {
	case sdk.FormatY16:
		return PixelFormatDepthU16, true
	case sdk.FormatRGB, sdk.FormatBGR:
		return PixelFormatRGB8, true
	case sdk.FormatRGBA, sdk.FormatBGRA:
		return PixelFormatRGBA8, true
	case sdk.FormatMJPG:
		return PixelFormatMJPEG, true
	case sdk.FormatYUYV, sdk.FormatNV12:
		return PixelFormatYUV422, true
	case sdk.FormatACCEL, sdk.FormatGYRO:
		return PixelFormatIMU, true
	default:
		return 0, false
	}
}
