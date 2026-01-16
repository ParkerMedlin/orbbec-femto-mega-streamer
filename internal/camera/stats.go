package camera

import "time"

// StreamStats captures runtime metrics for a single stream.
type StreamStats struct {
	Captured           uint64
	Dropped            uint64
	AvgLatency         time.Duration
	LastFrameID        uint64
	LastFrameTimestamp time.Time
}

// DeviceStats aggregates stats across all active streams on a device.
type DeviceStats struct {
	Streams map[StreamType]StreamStats
}

// Stats returns a snapshot of stream metrics.
func (s *Stream) Stats() StreamStats {
	if s == nil {
		return StreamStats{}
	}

	captured := s.captured.Load()
	totalLatency := s.totalLatency.Load()

	var avg time.Duration
	if captured > 0 {
		avg = time.Duration(totalLatency / int64(captured))
	}

	lastTs := s.lastTimestamp.Load()
	var ts time.Time
	if lastTs != 0 {
		ts = time.Unix(0, lastTs)
	}

	return StreamStats{
		Captured:           captured,
		Dropped:            s.dropped.Load(),
		AvgLatency:         avg,
		LastFrameID:        s.lastFrameID.Load(),
		LastFrameTimestamp: ts,
	}
}

// Stats returns stats for all streams on the device.
func (d *Device) Stats() DeviceStats {
	stats := DeviceStats{Streams: make(map[StreamType]StreamStats)}
	if d == nil {
		return stats
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	for t, s := range d.streams {
		stats.Streams[t] = s.Stats()
	}
	return stats
}
