package camera

import (
	"testing"
)

func TestFramePool(t *testing.T) {
	pool := newFramePool(3, 1024)

	f1 := pool.Get()
	f2 := pool.Get()
	f3 := pool.Get()

	if f1 == nil || f2 == nil || f3 == nil {
		t.Fatal("expected non-nil frames")
	}

	// Exhaust pool then release one and ensure reuse works.
	f1.Release()
	f4 := pool.Get()
	if f4 == nil {
		t.Fatal("expected frame after release")
	}

	// Ensure Release on nil-safe.
	var nilFrame *Frame
	nilFrame.Release()
}

func TestFrameClone(t *testing.T) {
	orig := &Frame{
		StreamType: StreamTypeColor,
		FrameID:    42,
		Timestamp:  1234,
		Width:      10,
		Height:     10,
		Format:     PixelFormatRGB8,
		Data:       []byte{1, 2, 3, 4},
	}
	clone := orig.Clone()

	if clone == orig || &clone.Data[0] == &orig.Data[0] {
		t.Fatal("expected deep copy")
	}
	clone.Data[0] = 9
	if orig.Data[0] == 9 {
		t.Fatal("mutating clone should not affect original")
	}
}

func BenchmarkFramePool(b *testing.B) {
	pool := newFramePool(10, 2048)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f := pool.Get()
			f.Data = append(f.Data[:0], make([]byte, 128)...)
			f.Release()
		}
	})
}
