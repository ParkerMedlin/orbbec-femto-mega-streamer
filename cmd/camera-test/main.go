package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera"
)

const targetFrames = 100

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	mgr, err := camera.NewDeviceManager(runCtx)
	if err != nil {
		log.Fatalf("init manager: %v", err)
	}
	defer mgr.Close()

	devices, err := mgr.Discover()
	if err != nil {
		log.Fatalf("discover devices: %v", err)
	}
	if len(devices) == 0 {
		log.Fatalf("no Orbbec devices detected")
	}

	dev, err := mgr.Open("")
	if err != nil {
		log.Fatalf("open device: %v", err)
	}
	defer dev.Close()

	depth := dev.Stream(camera.StreamTypeDepth)
	color := dev.Stream(camera.StreamTypeColor)

	if err := depth.Start(); err != nil {
		log.Fatalf("start depth stream: %v", err)
	}
	if err := color.Start(); err != nil {
		_ = depth.Stop()
		log.Fatalf("start color stream: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go collect(runCtx, &wg, "depth", depth)
	go collect(runCtx, &wg, "color", color)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("capture complete")
	case <-runCtx.Done():
		fmt.Printf("capture stopped early: %v\n", runCtx.Err())
	}

	_ = depth.Stop()
	_ = color.Stop()

	stats := dev.Stats()
	for streamType, s := range stats.Streams {
		fmt.Printf("%s stream: captured=%d dropped=%d avg_latency=%v last_frame_id=%d\n",
			streamType.String(), s.Captured, s.Dropped, s.AvgLatency, s.LastFrameID)
	}
}

func collect(ctx context.Context, wg *sync.WaitGroup, name string, stream *camera.Stream) {
	defer wg.Done()
	received := 0
	for received < targetFrames {
		select {
		case f := <-stream.Frames():
			if f != nil {
				received++
				f.Release()
			}
		case <-ctx.Done():
			return
		}
	}
	fmt.Printf("%s stream collected %d frames\n", name, received)
}
