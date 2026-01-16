# Orbbec Camera Integration Implementation Tasks

## Task Overview

This document breaks down the implementation of the Orbbec camera integration into actionable coding tasks. Each task is designed to be completed incrementally, with clear deliverables and requirements traceability.

**Total Tasks**: 20 tasks organized into 5 phases

**Requirements Reference**: This implementation addresses requirements from `requirements.md`

**Design Reference**: Technical approach defined in `design.md`

## Implementation Tasks

### Phase 1: Project Setup & CGo Foundation

- [ ] **1.1** Initialize Go module and project structure
  - **Description**: Create the base Go project with proper module initialization, directory structure, and build configuration for CGo
  - **Deliverables**:
    - `go.mod` with module name `github.com/[org]/orbbec-femto-mega-streamer`
    - Directory structure: `cmd/`, `internal/camera/`, `internal/camera/sdk/`
    - `Makefile` with build targets for Windows/Linux
  - **Requirements**: Technical Requirements - Project structure
  - **Dependencies**: None

- [ ] **1.2** Create CGo build configuration
  - **Description**: Set up CGo compiler flags pointing to the Orbbec SDK. Verify SDK headers are accessible and library can link
  - **Deliverables**:
    - `internal/camera/sdk/sdk.go` with CGo preamble (CFLAGS, LDFLAGS)
    - Build verification script that compiles successfully
  - **Requirements**: Technical Requirements - CGo bindings
  - **Dependencies**: 1.1
  - **Reference**: `orbbec-sdk/OrbbecSDK_v1.10.27/SDK/include/libobsensor/ObSensor.h`

- [ ] **1.3** Implement SDK error handling wrapper
  - **Description**: Create Go wrapper for `ob_error` handling. Implement `checkError()` function that extracts error details and frees C memory
  - **Deliverables**:
    - `internal/camera/sdk/error.go` - error checking and conversion
    - `internal/camera/errors.go` - Go error types (`SDKError`, `DeviceError`, `ConfigError`)
    - Unit tests for error type methods
  - **Requirements**: Error Handling Requirements
  - **Dependencies**: 1.2
  - **Reference**: `SDK/include/libobsensor/h/Error.h`

- [ ] **1.4** Define SDK type mappings
  - **Description**: Create Go constants and types that map to SDK enumerations (`ob_sensor_type`, `ob_format`, `ob_exception_type`)
  - **Deliverables**:
    - `internal/camera/sdk/types.go` - C enum mappings
    - `internal/camera/types.go` - Public Go types (`StreamType`, `PixelFormat`)
    - Unit tests verifying enum values match SDK
  - **Requirements**: Stream Capture Requirements - Pixel formats
  - **Dependencies**: 1.2
  - **Reference**: `SDK/include/libobsensor/h/ObTypes.h`

### Phase 2: SDK Wrapper Layer

- [ ] **2.1** Implement Context wrapper
  - **Description**: Wrap `ob_context` lifecycle (create, delete). Add logging configuration using `ob_set_logger_severity()`
  - **Deliverables**:
    - `internal/camera/sdk/context.go`
    - Functions: `CreateContext()`, `(*Context).Close()`, `(*Context).SetLogLevel()`
    - Integration test that creates and destroys context
  - **Requirements**: Device Connection Requirements - SDK initialization
  - **Dependencies**: 1.3
  - **Reference**: `SDK/include/libobsensor/h/Context.h`, `Example/c/Sample-HelloOrbbec`

- [ ] **2.2** Implement DeviceList wrapper
  - **Description**: Wrap `ob_device_list` operations for device enumeration. Extract device count and serial numbers
  - **Deliverables**:
    - `internal/camera/sdk/device_list.go`
    - Functions: `(*Context).QueryDeviceList()`, `(*DeviceList).Count()`, `(*DeviceList).SerialNumber(index)`, `(*DeviceList).Close()`
  - **Requirements**: Device Connection Requirements - Enumerate devices
  - **Dependencies**: 2.1
  - **Reference**: `SDK/include/libobsensor/h/Device.h`

- [ ] **2.3** Implement Device wrapper
  - **Description**: Wrap `ob_device` operations. Get device info (name, serial, firmware version, connection type)
  - **Deliverables**:
    - `internal/camera/sdk/device.go`
    - Functions: `(*DeviceList).GetDevice(index)`, `(*Device).GetInfo()`, `(*Device).Close()`
    - Struct: `DeviceInfo` with all metadata fields
  - **Requirements**: Device Connection Requirements - Device metadata
  - **Dependencies**: 2.2
  - **Reference**: `SDK/include/libobsensor/h/Device.h`

- [ ] **2.4** Implement Pipeline wrapper
  - **Description**: Wrap `ob_pipeline` for high-level streaming. Handle config creation, stream enable, start/stop
  - **Deliverables**:
    - `internal/camera/sdk/pipeline.go`
    - Functions: `CreatePipelineWithDevice()`, `(*Pipeline).EnableStream()`, `(*Pipeline).Start()`, `(*Pipeline).Stop()`, `(*Pipeline).Close()`
  - **Requirements**: Stream Capture Requirements
  - **Dependencies**: 2.3
  - **Reference**: `SDK/include/libobsensor/h/Pipeline.h`, `Example/c/Sample-DepthViewer`

- [ ] **2.5** Implement Frame and FrameSet wrappers
  - **Description**: Wrap `ob_frame` and `ob_frameset` for frame data access. Extract timestamps, dimensions, format, raw data pointer
  - **Deliverables**:
    - `internal/camera/sdk/frame.go`
    - Functions: `(*Pipeline).WaitForFrameSet()`, `(*FrameSet).GetDepthFrame()`, `(*FrameSet).GetColorFrame()`, `(*Frame).GetData()`, `(*Frame).GetTimestamp()`, etc.
  - **Requirements**: Frame Data Requirements
  - **Dependencies**: 2.4
  - **Reference**: `SDK/include/libobsensor/h/Frame.h`

- [ ] **2.6** Implement StreamProfile wrapper
  - **Description**: Wrap `ob_stream_profile` for querying available stream configurations (resolutions, FPS, formats)
  - **Deliverables**:
    - `internal/camera/sdk/stream_profile.go`
    - Functions: `(*Device).GetStreamProfileList()`, `(*StreamProfile).GetWidth()`, `(*StreamProfile).GetHeight()`, `(*StreamProfile).GetFPS()`, `(*StreamProfile).GetFormat()`
  - **Requirements**: Stream Configuration Requirements
  - **Dependencies**: 2.3
  - **Reference**: `SDK/include/libobsensor/h/StreamProfile.h`

### Phase 3: Go API Layer

- [ ] **3.1** Implement DeviceManager
  - **Description**: Create high-level Go API for device management. Handle device discovery, opening, and lifecycle
  - **Deliverables**:
    - `internal/camera/manager.go`
    - Type: `DeviceManager` struct
    - Functions: `NewDeviceManager()`, `(*DeviceManager).Discover()`, `(*DeviceManager).Open(serial)`, `(*DeviceManager).Close()`
    - Unit tests with mock SDK
  - **Requirements**: Device Connection Requirements
  - **Dependencies**: 2.3

- [ ] **3.2** Implement Device hot-plug callbacks
  - **Description**: Register for device add/remove events using `ob_set_device_changed_callback()`. Export CGo callback, route events to Go channels
  - **Deliverables**:
    - `internal/camera/sdk/callbacks.go` - CGo callback exports
    - `internal/camera/manager.go` - event handling methods
    - Functions: `(*DeviceManager).OnDeviceAdded()`, `(*DeviceManager).OnDeviceRemoved()` (return channels)
  - **Requirements**: Error Handling Requirements - Camera disconnect/reconnect
  - **Dependencies**: 3.1
  - **Reference**: `Example/c/Sample-HotPlugin`

- [ ] **3.3** Implement StreamConfig and validation
  - **Description**: Create configuration struct for stream settings. Validate against supported profiles from device
  - **Deliverables**:
    - `internal/camera/config.go`
    - Type: `StreamConfig` struct
    - Functions: `DefaultDepthConfig()`, `DefaultColorConfig()`, `(*StreamConfig).Validate(device)`
    - Unit tests for validation logic
  - **Requirements**: Stream Configuration Requirements
  - **Dependencies**: 2.6

- [ ] **3.4** Implement Frame type and buffer pool
  - **Description**: Create Go Frame struct with pooled buffer management to minimize allocations
  - **Deliverables**:
    - `internal/camera/frame.go`
    - Types: `Frame`, `framePool`
    - Functions: `newFramePool()`, `(*framePool).Get()`, `(*Frame).Release()`, `(*Frame).Clone()`
    - Benchmark tests for pool performance
  - **Requirements**: Frame Data Requirements, Performance Requirements
  - **Dependencies**: 1.4

### Phase 4: Frame Capture Pipeline

- [ ] **4.1** Implement Stream type with capture goroutine
  - **Description**: Create Stream type that runs capture loop in dedicated goroutine, copies frames to pool buffers, sends on channel
  - **Deliverables**:
    - `internal/camera/stream.go`
    - Type: `Stream` struct
    - Functions: `(*Device).NewStream()`, `(*Stream).Start()`, `(*Stream).Stop()`, `(*Stream).Frames()` (returns channel)
  - **Requirements**: Stream Capture Requirements
  - **Dependencies**: 2.5, 3.4

- [ ] **4.2** Implement runtime stream reconfiguration
  - **Description**: Allow changing stream resolution/FPS/format while running. Handle pipeline restart internally
  - **Deliverables**:
    - `internal/camera/stream.go` - `(*Stream).SetConfig()`
    - Proper synchronization during reconfiguration
    - Tests verifying config changes apply correctly
  - **Requirements**: Stream Configuration Requirements - Apply within 500ms
  - **Dependencies**: 4.1

- [ ] **4.3** Implement Device type with multi-stream support
  - **Description**: Device holds multiple Stream instances (depth, color, IR, IMU). Coordinate lifecycle
  - **Deliverables**:
    - `internal/camera/device.go`
    - Type: `Device` struct
    - Functions: `(*Device).Stream(StreamType)`, `(*Device).Info()`, `(*Device).Close()`
  - **Requirements**: Stream Capture Requirements - Multiple stream types
  - **Dependencies**: 4.1

### Phase 5: Integration & Testing

- [ ] **5.1** Create integration test harness
  - **Description**: Build test infrastructure for real hardware testing. Skip tests gracefully when no camera connected
  - **Deliverables**:
    - `internal/camera/integration_test.go` (build tag: integration)
    - Test: device discovery, connection, basic frame capture
    - CI configuration to skip integration tests without hardware
  - **Requirements**: Testing Requirements
  - **Dependencies**: 4.3

- [ ] **5.2** Implement capture statistics
  - **Description**: Track frames captured, dropped, latency metrics. Expose via device/stream methods
  - **Deliverables**:
    - `internal/camera/stats.go`
    - Type: `StreamStats` (frames captured, dropped, avg latency)
    - Functions: `(*Stream).Stats()`, `(*Device).Stats()`
  - **Requirements**: Performance Requirements
  - **Dependencies**: 4.1

- [ ] **5.3** Create example application
  - **Description**: Build minimal example demonstrating camera connection, stream capture, and frame access
  - **Deliverables**:
    - `cmd/camera-test/main.go`
    - Example: connect, enable depth+color, capture 100 frames, print stats
  - **Requirements**: All requirements (end-to-end validation)
  - **Dependencies**: 5.2

- [ ] **5.4** Add graceful shutdown handling
  - **Description**: Ensure proper cleanup on SIGINT/SIGTERM. Stop streams, release SDK resources, avoid crashes
  - **Deliverables**:
    - Signal handling in example app
    - `(*DeviceManager).Close()` properly stops all streams
    - Test: verify no resource leaks on shutdown
  - **Requirements**: Technical Requirements - Graceful shutdown
  - **Dependencies**: 5.3

## Task Guidelines

### Task Completion Criteria
Each task is considered complete when:
- [ ] All deliverables are implemented and functional
- [ ] Unit tests are written and passing (where applicable)
- [ ] Code compiles without warnings
- [ ] CGo memory management is correct (no leaks)
- [ ] Error handling covers all SDK failure cases

### Task Dependencies
- Tasks should be completed in order within each phase
- Phase 2 depends on Phase 1 completion
- Phase 3 can partially overlap with Phase 2
- Phase 4 requires Phase 2 and 3
- Phase 5 requires Phase 4

### Testing Requirements
- **Unit Tests**: Required for all Go logic (config validation, error types, frame pool)
- **Integration Tests**: Required for SDK wrapper functions (with real hardware)
- **Benchmark Tests**: Required for frame pool and capture loop performance

### Code Quality Standards
- All CGo calls must check and handle errors
- All C memory allocations must have corresponding frees
- Use `defer` for cleanup where appropriate
- Document all exported types and functions
- Use Go naming conventions (no C-style names in public API)

## Resource Requirements

### Development Environment
- Go 1.21+ with CGo enabled
- GCC/MinGW (Windows) or GCC (Linux) for C compilation
- Orbbec SDK v1.10.27 (included in `orbbec-sdk/`)
- Orbbec Femto Mega camera (for integration testing)

### External Dependencies
- Orbbec SDK libraries must be in library path at runtime:
  - Windows: `OrbbecSDK.dll`, `live555.dll`, `ob_usb.dll`, `depthengine_2_0.dll`
  - Linux: `libOrbbecSDK.so`

### Build Commands
```bash
# Windows (PowerShell)
$env:CGO_ENABLED=1
go build ./...

# Linux
CGO_ENABLED=1 go build ./...

# Run tests (unit only)
go test ./internal/camera/...

# Run integration tests (requires camera)
go test -tags=integration ./internal/camera/...
```

## Risk Mitigation

### Technical Risks
- **Risk**: CGo callback threading issues (SDK callbacks from non-Go threads)
  - **Mitigation**: Use `runtime.LockOSThread()` in callback handlers, channel-based event passing
  - **Affected Tasks**: 3.2

- **Risk**: Memory corruption from incorrect buffer handling
  - **Mitigation**: Always copy frame data immediately, never hold C pointers across calls
  - **Affected Tasks**: 2.5, 3.4, 4.1

- **Risk**: SDK version incompatibility
  - **Mitigation**: Pin to v1.10.27, test with exact version, document version requirements
  - **Affected Tasks**: All Phase 2 tasks

---

**Task Status**: Not Started

**Current Phase**: Phase 1

**Overall Progress**: 0/20 tasks completed (0%)

**Last Updated**: 2025-01-15
