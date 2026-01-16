package camera

import "testing"

type fakeProfileProvider struct {
	profiles map[StreamType][]StreamProfileInfo
	err      error
}

func (f *fakeProfileProvider) Profiles(t StreamType) ([]StreamProfileInfo, error) {
	return f.profiles[t], f.err
}

func TestStreamConfigValidate(t *testing.T) {
	provider := &fakeProfileProvider{
		profiles: map[StreamType][]StreamProfileInfo{
			StreamTypeDepth: {
				{Width: 640, Height: 576, FPS: 30, Format: PixelFormatDepthU16},
			},
		},
	}

	tests := []struct {
		name    string
		cfg     StreamConfig
		wantErr bool
	}{
		{"valid depth", StreamConfig{Type: StreamTypeDepth, Enabled: true, Width: 640, Height: 576, FPS: 30, Format: PixelFormatDepthU16}, false},
		{"invalid resolution", StreamConfig{Type: StreamTypeDepth, Enabled: true, Width: 999, Height: 576, FPS: 30, Format: PixelFormatDepthU16}, true},
		{"invalid fps", StreamConfig{Type: StreamTypeDepth, Enabled: true, Width: 640, Height: 576, FPS: 0, Format: PixelFormatDepthU16}, true},
		{"unsupported profile", StreamConfig{Type: StreamTypeDepth, Enabled: true, Width: 320, Height: 240, FPS: 30, Format: PixelFormatDepthU16}, true},
		{"disabled ok", StreamConfig{Type: StreamTypeDepth, Enabled: false, Width: 0, Height: 0, FPS: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate(provider)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfigs(t *testing.T) {
	depth := DefaultDepthConfig()
	if depth.Format != PixelFormatDepthU16 || !depth.Enabled {
		t.Fatalf("unexpected depth default: %+v", depth)
	}

	color := DefaultColorConfig()
	if color.Format != PixelFormatRGB8 || color.Width == 0 || !color.Enabled {
		t.Fatalf("unexpected color default: %+v", color)
	}
}
