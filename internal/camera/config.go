package camera

import "fmt"

// StreamProfileInfo describes a supported stream profile exposed by the device.
type StreamProfileInfo struct {
	Width  uint16
	Height uint16
	FPS    uint8
	Format PixelFormat
}

// ProfileProvider exposes supported stream profiles for validation.
type ProfileProvider interface {
	Profiles(streamType StreamType) ([]StreamProfileInfo, error)
}

// StreamConfig defines the configuration for a single stream.
type StreamConfig struct {
	Type    StreamType
	Enabled bool
	Width   uint16
	Height  uint16
	FPS     uint8
	Format  PixelFormat
}

// Validate checks the config for basic correctness and compatibility with the provider.
func (c *StreamConfig) Validate(provider ProfileProvider) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	if !c.Enabled {
		return nil
	}

	if c.Width == 0 || c.Height == 0 {
		return &ConfigError{Stream: c.Type, Field: "resolution", Value: fmt.Sprintf("%dx%d", c.Width, c.Height), Reason: "must be > 0"}
	}
	if c.FPS == 0 {
		return &ConfigError{Stream: c.Type, Field: "fps", Value: c.FPS, Reason: "must be > 0"}
	}

	if provider == nullProvider {
		return nil
	}

	if provider != nil {
		profiles, err := provider.Profiles(c.Type)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			return &ConfigError{Stream: c.Type, Field: "profile", Value: "none", Reason: "device reports no profiles"}
		}
		for _, p := range profiles {
			if p.Width == c.Width && p.Height == c.Height && p.FPS == c.FPS && p.Format == c.Format {
				return nil
			}
		}
		return &ConfigError{Stream: c.Type, Field: "profile", Value: fmt.Sprintf("%dx%d@%d %s", c.Width, c.Height, c.FPS, c.Format), Reason: "unsupported by device"}
	}

	return nil
}

// DefaultDepthConfig returns a sensible depth stream default.
func DefaultDepthConfig() StreamConfig {
	return StreamConfig{
		Type:    StreamTypeDepth,
		Enabled: true,
		Width:   640,
		Height:  576,
		FPS:     30,
		Format:  PixelFormatDepthU16,
	}
}

// DefaultColorConfig returns a sensible color stream default.
func DefaultColorConfig() StreamConfig {
	return StreamConfig{
		Type:    StreamTypeColor,
		Enabled: true,
		Width:   1920,
		Height:  1080,
		FPS:     30,
		Format:  PixelFormatRGB8,
	}
}

// nullProvider skips device profile validation while keeping nil distinct.
var nullProvider ProfileProvider
