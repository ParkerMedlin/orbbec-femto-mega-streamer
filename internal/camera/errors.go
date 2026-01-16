package camera

import "fmt"

// DeviceError wraps errors tied to a specific device operation.
type DeviceError struct {
	Serial string
	Op     string
	Err    error
}

func (e *DeviceError) Error() string {
	return fmt.Sprintf("device %s %s failed: %v", e.Serial, e.Op, e.Err)
}

func (e *DeviceError) Unwrap() error { return e.Err }

// ConfigError indicates invalid configuration values for a stream.
type ConfigError struct {
	Stream StreamType
	Field  string
	Value  interface{}
	Reason string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("invalid %s config: %s=%v (%s)", e.Stream, e.Field, e.Value, e.Reason)
}
