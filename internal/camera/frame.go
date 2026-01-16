package camera

// Frame represents a single captured frame from a stream.
type Frame struct {
	StreamType StreamType
	FrameID    uint64
	Timestamp  int64
	Width      uint16
	Height     uint16
	Format     PixelFormat
	Data       []byte
	pool       *framePool
}

// Size returns the frame payload size in bytes.
func (f *Frame) Size() int {
	if f == nil {
		return 0
	}
	return len(f.Data)
}

// Release returns the frame buffer to its pool if available.
func (f *Frame) Release() {
	if f == nil || f.pool == nil {
		return
	}
	f.Data = f.Data[:0]
	f.pool.Put(f)
}

// Clone makes a deep copy of the frame data with no pool association.
func (f *Frame) Clone() *Frame {
	if f == nil {
		return nil
	}
	clone := *f
	clone.pool = nil
	if f.Data != nil {
		clone.Data = append([]byte(nil), f.Data...)
	}
	return &clone
}

// framePool manages preallocated frame buffers to reduce GC churn.
type framePool struct {
	buffers chan *Frame
	size    int
}

func newFramePool(count, bufferSize int) *framePool {
	if count <= 0 {
		count = 1
	}
	p := &framePool{
		buffers: make(chan *Frame, count),
		size:    bufferSize,
	}
	for i := 0; i < count; i++ {
		p.buffers <- &Frame{
			Data: make([]byte, 0, bufferSize),
			pool: p,
		}
	}
	return p
}

func (p *framePool) Get() *Frame {
	if p == nil {
		return nil
	}
	select {
	case f := <-p.buffers:
		return f
	default:
		return &Frame{
			Data: make([]byte, 0, p.size),
			pool: p,
		}
	}
}

func (p *framePool) Put(f *Frame) {
	if p == nil || f == nil {
		return
	}
	f.pool = p
	f.Data = f.Data[:0]
	select {
	case p.buffers <- f:
	default:
		// Pool full; drop extra buffer to let GC reclaim.
	}
}
