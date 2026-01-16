//go:build integration
// +build integration

package camera

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// TestIntegrationCapture exercises a basic depth stream capture against real hardware.
func TestIntegrationCapture(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("integration test requires Windows SDK and camera hardware")
	}

	mgr, err := NewDeviceManager(context.Background())
	if err != nil {
		t.Skipf("unable to init SDK: %v", err)
	}
	defer mgr.Close()

	devices, err := mgr.Discover()
	if err != nil {
		t.Skipf("device discovery failed: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("no cameras connected")
	}

	dev, err := mgr.Open(devices[0].Serial)
	if err != nil {
		t.Fatalf("open device: %v", err)
	}
	defer dev.Close()

	stream := dev.Stream(StreamTypeDepth)
	if stream == nil {
		t.Fatalf("depth stream is nil")
	}
	if err := stream.Start(); err != nil {
		t.Fatalf("start depth stream: %v", err)
	}
	defer stream.Stop()

	received := 0
	timeout := time.After(5 * time.Second)
	for received < 5 {
		select {
		case f := <-stream.Frames():
			if f != nil {
				received++
				f.Release()
			}
		case <-timeout:
			t.Fatalf("timed out waiting for frames; received %d", received)
		}
	}

	stats := stream.Stats()
	if stats.Captured == 0 {
		t.Fatalf("expected stats to record captured frames")
	}
}
