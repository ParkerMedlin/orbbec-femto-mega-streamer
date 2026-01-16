//go:build !windows || !cgo

package camera

import (
	"fmt"

	"github.com/ParkerMedlin/orbbec-femto-mega-streamer/internal/camera/sdk"
)

func init() {
	pipelineFactory = func(dev *sdk.Device) (streamPipeline, error) {
		return nil, fmt.Errorf("orbbec sdk pipeline requires Windows with cgo enabled")
	}
}
