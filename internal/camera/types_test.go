package camera

import "testing"

func TestPixelFormatBytesPerPixel(t *testing.T) {
	tests := []struct {
		fmt  PixelFormat
		want int
	}{
		{PixelFormatDepthU16, 2},
		{PixelFormatDepthF32, 4},
		{PixelFormatRGB8, 3},
		{PixelFormatRGBA8, 4},
		{PixelFormatMJPEG, 0},
	}

	for _, tt := range tests {
		if got := tt.fmt.BytesPerPixel(); got != tt.want {
			t.Fatalf("%v bytes-per-pixel = %d, want %d", tt.fmt, got, tt.want)
		}
	}
}
