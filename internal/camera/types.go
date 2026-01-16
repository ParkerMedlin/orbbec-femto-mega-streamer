package camera

// StreamType identifies the type of data stream.
type StreamType uint8

const (
	StreamTypeDepth StreamType = iota
	StreamTypeColor
	StreamTypeIR
	StreamTypeIMU
)

func (s StreamType) String() string {
	switch s {
	case StreamTypeDepth:
		return "depth"
	case StreamTypeColor:
		return "color"
	case StreamTypeIR:
		return "ir"
	case StreamTypeIMU:
		return "imu"
	default:
		return "unknown"
	}
}

// PixelFormat identifies the pixel data format.
type PixelFormat uint8

const (
	PixelFormatDepthU16 PixelFormat = iota // uint16 depth in millimeters
	PixelFormatDepthF32                    // float32 depth in meters
	PixelFormatRGB8                        // 3 bytes per pixel
	PixelFormatRGBA8                       // 4 bytes per pixel
	PixelFormatMJPEG                       // compressed JPEG
	PixelFormatYUV422                      // packed YUV 4:2:2
	PixelFormatIR16                        // uint16 infrared intensity
	PixelFormatIMU                         // IMU sample struct
)

// BytesPerPixel returns the nominal bytes per pixel for uncompressed formats.
// For compressed formats (e.g., MJPEG) it returns 0 because size is variable.
func (p PixelFormat) BytesPerPixel() int {
	switch p {
	case PixelFormatDepthU16, PixelFormatIR16:
		return 2
	case PixelFormatDepthF32:
		return 4
	case PixelFormatRGB8:
		return 3
	case PixelFormatRGBA8, PixelFormatYUV422:
		return 4
	case PixelFormatMJPEG, PixelFormatIMU:
		return 0
	default:
		return 0
	}
}

func (p PixelFormat) String() string {
	switch p {
	case PixelFormatDepthU16:
		return "depth_u16"
	case PixelFormatDepthF32:
		return "depth_f32"
	case PixelFormatRGB8:
		return "rgb8"
	case PixelFormatRGBA8:
		return "rgba8"
	case PixelFormatMJPEG:
		return "mjpeg"
	case PixelFormatYUV422:
		return "yuv422"
	case PixelFormatIR16:
		return "ir16"
	case PixelFormatIMU:
		return "imu"
	default:
		return "unknown"
	}
}
