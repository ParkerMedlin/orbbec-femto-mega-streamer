package camera

import (
	"errors"
	"testing"
)

func TestDeviceError(t *testing.T) {
	base := errors.New("boom")
	err := &DeviceError{Serial: "1234", Op: "connect", Err: base}

	if got := err.Error(); got == "" || got == "boom" {
		t.Fatalf("unexpected error string: %q", got)
	}

	if !errors.Is(err, base) {
		t.Fatalf("expected errors.Is to unwrap base error")
	}
}

func TestConfigError(t *testing.T) {
	err := &ConfigError{Stream: StreamTypeDepth, Field: "fps", Value: 999, Reason: "unsupported"}
	got := err.Error()
	if got == "" || got == "unsupported" {
		t.Fatalf("unexpected error string: %q", got)
	}
}
