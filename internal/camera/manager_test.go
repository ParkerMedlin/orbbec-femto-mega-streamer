package camera

import (
	"context"
	"errors"
	"testing"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

type fakeAdapter struct {
	devices []DeviceInfo
	cb      sdk.DeviceChangedCallback
	closed  bool
}

func (f *fakeAdapter) ListDevices() ([]DeviceInfo, error) {
	return f.devices, nil
}

func (f *fakeAdapter) OpenDevice(serial string) (*sdk.Device, *DeviceInfo, error) {
	if len(f.devices) == 0 {
		return nil, nil, errors.New("no devices")
	}
	if serial == "" {
		serial = f.devices[0].Serial
	}
	for _, d := range f.devices {
		if d.Serial == serial {
			info := d
			return &sdk.Device{}, &info, nil
		}
	}
	return nil, nil, errors.New("not found")
}

func (f *fakeAdapter) SetDeviceChangedCallback(cb sdk.DeviceChangedCallback) (func(), error) {
	f.cb = cb
	return func() { f.cb = nil }, nil
}

func (f *fakeAdapter) Close() error {
	f.closed = true
	return nil
}

func (f *fakeAdapter) emit(added, removed []string) {
	if f.cb != nil {
		f.cb(added, removed)
	}
}

func withFakeAdapter(t *testing.T, f *fakeAdapter) func() {
	t.Helper()
	orig := newSDKAdapter
	newSDKAdapter = func() (sdkAdapter, error) {
		return f, nil
	}
	return func() { newSDKAdapter = orig }
}

func TestDeviceManagerDiscover(t *testing.T) {
	fake := &fakeAdapter{
		devices: []DeviceInfo{{Serial: "123", Name: "cam"}},
	}
	restore := withFakeAdapter(t, fake)
	defer restore()

	mgr, err := NewDeviceManager(context.Background())
	if err != nil {
		t.Fatalf("NewDeviceManager error: %v", err)
	}
	defer mgr.Close()

	devs, err := mgr.Discover()
	if err != nil {
		t.Fatalf("Discover error: %v", err)
	}
	if len(devs) != 1 || devs[0].Serial != "123" {
		t.Fatalf("unexpected devices: %+v", devs)
	}
}

func TestDeviceManagerOpenDefault(t *testing.T) {
	fake := &fakeAdapter{
		devices: []DeviceInfo{
			{Serial: "A"},
			{Serial: "B"},
		},
	}
	restore := withFakeAdapter(t, fake)
	defer restore()

	mgr, err := NewDeviceManager(context.Background())
	if err != nil {
		t.Fatalf("NewDeviceManager error: %v", err)
	}
	defer mgr.Close()

	dev, err := mgr.Open("")
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	if dev.Info().Serial != "A" {
		t.Fatalf("expected first device serial A, got %s", dev.Info().Serial)
	}

	// Open again should reuse.
	dev2, err := mgr.Open("A")
	if err != nil {
		t.Fatalf("Open again error: %v", err)
	}
	if dev != dev2 {
		t.Fatal("expected same device instance to be reused")
	}
}

func TestDeviceManagerHotplug(t *testing.T) {
	fake := &fakeAdapter{
		devices: []DeviceInfo{{Serial: "A"}},
	}
	restore := withFakeAdapter(t, fake)
	defer restore()

	mgr, err := NewDeviceManager(context.Background())
	if err != nil {
		t.Fatalf("NewDeviceManager error: %v", err)
	}
	defer mgr.Close()

	fake.emit([]string{"B"}, []string{"A"})

	select {
	case s := <-mgr.OnDeviceAdded():
		if s != "B" {
			t.Fatalf("unexpected added serial %s", s)
		}
	default:
		t.Fatal("expected added event")
	}

	select {
	case s := <-mgr.OnDeviceRemoved():
		if s != "A" {
			t.Fatalf("unexpected removed serial %s", s)
		}
	default:
		t.Fatal("expected removed event")
	}
}
