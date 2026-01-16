//go:build windows && cgo

package sdk

import "testing"

func TestFrameFormatValues(t *testing.T) {
	if FormatRGB != FrameFormat(C.OB_FORMAT_RGB) {
		t.Fatalf("FormatRGB mismatch: got %d want %d", FormatRGB, FrameFormat(C.OB_FORMAT_RGB))
	}
	if FormatY16 != FrameFormat(C.OB_FORMAT_Y16) {
		t.Fatalf("FormatY16 mismatch: got %d want %d", FormatY16, FrameFormat(C.OB_FORMAT_Y16))
	}
}

func TestSensorTypeValues(t *testing.T) {
	if SensorDepth != SensorType(C.OB_SENSOR_DEPTH) {
		t.Fatalf("SensorDepth mismatch")
	}
	if SensorColor != SensorType(C.OB_SENSOR_COLOR) {
		t.Fatalf("SensorColor mismatch")
	}
}
