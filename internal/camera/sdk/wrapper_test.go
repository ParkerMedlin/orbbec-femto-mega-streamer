//go:build windows && cgo

package sdk

import "testing"

// These tests are lightweight compile/guard checks. They do not require hardware.

func TestContextCloseNil(t *testing.T) {
	var ctx *Context
	if err := ctx.Close(); err != nil {
		t.Fatalf("expected nil error for nil context close: %v", err)
	}
}

func TestDeviceListCloseNil(t *testing.T) {
	var dl *DeviceList
	if err := dl.Close(); err != nil {
		t.Fatalf("expected nil error for nil device list close: %v", err)
	}
}

func TestPipelineCloseNil(t *testing.T) {
	var p *Pipeline
	if err := p.Close(); err != nil {
		t.Fatalf("expected nil error for nil pipeline close: %v", err)
	}
}

func TestFrameCloseNil(t *testing.T) {
	var f *Frame
	if err := f.Close(); err != nil {
		t.Fatalf("expected nil error for nil frame close: %v", err)
	}
}

func TestFrameSetCloseNil(t *testing.T) {
	var fs *FrameSet
	if err := fs.Close(); err != nil {
		t.Fatalf("expected nil error for nil frameset close: %v", err)
	}
}

func TestStreamProfileCloseNil(t *testing.T) {
	var sp *StreamProfile
	if err := sp.Close(); err != nil {
		t.Fatalf("expected nil error for nil stream profile close: %v", err)
	}
}

func TestStreamProfileListCloseNil(t *testing.T) {
	var spl *StreamProfileList
	if err := spl.Close(); err != nil {
		t.Fatalf("expected nil error for nil stream profile list close: %v", err)
	}
}
