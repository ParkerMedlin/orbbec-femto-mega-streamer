# Orbbec Camera Integration Requirements

## 1. Introduction

This document specifies the requirements for integrating the Orbbec Femto Mega depth camera into the Go-based network streaming application. The integration provides the foundation for capturing depth, color, IR, and IMU data streams that will be transmitted over the network.

**Architecture Overview**: The application uses CGo bindings to interface with the Orbbec SDK (C/C++ API). Each stream type runs in its own goroutine, capturing frames from the camera and passing them to the streaming pipeline. The integration must support runtime configuration changes and handle camera disconnect/reconnect scenarios gracefully.

## 2. User Stories

### Camera Operator
- **As a camera operator**, I want to connect to an Orbbec Femto Mega camera, so that I can capture depth and color data
- **As a camera operator**, I want to enable/disable individual streams at runtime, so that I can conserve bandwidth and processing when certain data isn't needed
- **As a camera operator**, I want to configure stream resolution and frame rate, so that I can balance quality vs performance for my use case

### System Integrator
- **As a system integrator**, I want the application to auto-detect connected cameras, so that I don't need to manually specify device paths
- **As a system integrator**, I want to specify startup defaults via configuration file, so that the system starts with appropriate settings
- **As a system integrator**, I want detailed logging of camera operations, so that I can diagnose issues in production

### Application Developer
- **As an application developer**, I want a clean Go interface for camera operations, so that I can easily integrate camera data into the streaming pipeline
- **As an application developer**, I want frame data with metadata (timestamp, resolution, format), so that I can properly process and transmit frames
- **As an application developer**, I want notification of camera events (connect/disconnect), so that I can handle hardware changes gracefully

## 3. Acceptance Criteria

### Device Connection Requirements
- **WHEN** the application starts, **THEN** the system **SHALL** enumerate all connected Orbbec devices
- **WHEN** a valid Orbbec Femto Mega is detected, **THEN** the system **SHALL** establish a connection within 5 seconds
- **WHEN** no camera is detected, **THEN** the system **SHALL** log an error and continue running (waiting for camera connection)
- **IF** multiple cameras are connected, **THEN** the system **SHALL** use the first detected camera (Phase 1) or allow selection by serial number (Phase 2)

### Stream Capture Requirements
- **WHEN** a stream is enabled, **THEN** the system **SHALL** begin capturing frames at the configured resolution and FPS
- **WHEN** a depth stream is active, **THEN** the system **SHALL** provide 16-bit depth values in millimeters
- **WHEN** a color stream is active, **THEN** the system **SHALL** provide RGB/RGBA/MJPEG frames as configured
- **WHEN** an IR stream is active, **THEN** the system **SHALL** provide infrared intensity data
- **WHEN** IMU stream is enabled, **THEN** the system **SHALL** provide 6DoF accelerometer and gyroscope data

### Stream Configuration Requirements
- **WHEN** resolution is changed, **THEN** the system **SHALL** apply the new resolution within 500ms
- **WHEN** FPS is changed, **THEN** the system **SHALL** apply the new frame rate within 500ms
- **WHEN** pixel format is changed, **THEN** the system **SHALL** apply the new format within 500ms
- **IF** an invalid configuration is requested, **THEN** the system **SHALL** reject it and maintain current settings
- **WHEN** a stream is disabled, **THEN** the system **SHALL** stop capturing and free associated resources

### Frame Data Requirements
- **WHEN** a frame is captured, **THEN** the system **SHALL** include timestamp (nanosecond precision)
- **WHEN** a frame is captured, **THEN** the system **SHALL** include frame dimensions (width, height)
- **WHEN** a frame is captured, **THEN** the system **SHALL** include pixel format identifier
- **WHEN** a frame is captured, **THEN** the system **SHALL** include raw pixel data buffer

### Error Handling Requirements
- **WHEN** camera disconnects unexpectedly, **THEN** the system **SHALL** detect within 1 second and emit a disconnect event
- **WHEN** camera reconnects, **THEN** the system **SHALL** auto-reconnect and resume previous configuration
- **IF** frame capture fails, **THEN** the system **SHALL** log the error and continue attempting capture
- **WHEN** SDK returns an error, **THEN** the system **SHALL** translate it to a meaningful Go error

## 4. Technical Architecture

### CGo Integration
- **Binding Approach**: Direct CGo bindings to Orbbec SDK C API
- **Memory Management**: Go-managed buffers with careful C memory handling
- **Thread Safety**: SDK calls coordinated through Go channels

### Go Package Structure
- **Package**: `internal/camera`
- **Key Types**: `Device`, `Stream`, `Frame`, `Config`
- **Interfaces**: `FrameHandler` for downstream consumers

### Key Libraries & Dependencies
- **Orbbec SDK**: Official C/C++ SDK for Femto Mega (v1.x or v2.x)
- **CGo**: Go's C interoperability layer
- **Configuration**: YAML/JSON config file support

## 5. Feature Specifications

### Core Features
1. **Device Discovery**: Enumerate and connect to Orbbec Femto Mega cameras
2. **Depth Stream Capture**: ToF sensor data at configurable resolution/FPS
3. **Color Stream Capture**: RGB video at configurable resolution/FPS/format
4. **Stream Configuration**: Runtime enable/disable, resolution, FPS, format changes
5. **Frame Pipeline**: Channel-based frame delivery to consumers

### Supported Configurations

#### Depth Stream
| Resolution | FPS Options | Pixel Formats |
|------------|-------------|---------------|
| 640x576    | 15, 30      | uint16 (mm)   |
| 512x512    | 15, 30      | uint16 (mm)   |
| 320x288    | 15, 30      | float32 (m)   |

#### Color Stream
| Resolution | FPS Options | Pixel Formats |
|------------|-------------|---------------|
| 3840x2160  | 15, 30      | MJPEG         |
| 1920x1080  | 15, 30      | RGB, MJPEG    |
| 1280x720   | 15, 30      | RGB, RGBA     |
| 640x480    | 15, 30, 60  | RGB, RGBA     |

#### IR Stream (Optional)
| Resolution | FPS Options | Pixel Formats |
|------------|-------------|---------------|
| 640x576    | 15, 30      | uint16        |

#### IMU Stream (Optional)
| Data Type | Rate | Format |
|-----------|------|--------|
| 6DoF      | 200Hz | Accel XYZ + Gyro XYZ |

## 6. Success Criteria

### Technical Performance
- **WHEN** streaming at 30 FPS, **THEN** frame latency **SHALL** be less than 50ms from capture to delivery
- **WHEN** streaming at 30 FPS, **THEN** CPU usage **SHALL** be less than 20% per stream on target hardware
- **WHEN** running for 24 hours, **THEN** memory usage **SHALL** remain stable (no leaks)

### Reliability
- **WHEN** tested over 1000 connect/disconnect cycles, **THEN** success rate **SHALL** exceed 99%
- **WHEN** camera firmware is compatible, **THEN** connection success rate **SHALL** exceed 99.9%

## 7. Assumptions and Dependencies

### Technical Assumptions
- Orbbec SDK is installed and accessible on the build system
- Target platform supports USB 3.0 for full bandwidth
- Go version 1.21+ with CGo support enabled
- GCC/Clang available for CGo compilation

### External Dependencies
- Orbbec SDK (libOrbbecSDK)
- USB drivers for Orbbec cameras
- libudev (Linux) or equivalent platform USB stack

## 8. Constraints and Limitations

### Technical Constraints
- Maximum one active depth stream per camera (hardware limitation)
- Color and depth streams may have alignment requirements
- USB bandwidth limits simultaneous high-resolution streams
- CGo introduces some overhead vs pure C

### Platform Constraints
- Initial support: Linux (primary), Windows (secondary)
- macOS support dependent on Orbbec SDK availability

## 9. Risk Assessment

### Technical Risks
- **Risk**: Orbbec SDK version incompatibility
  - **Likelihood**: Medium
  - **Impact**: High
  - **Mitigation**: Pin to specific SDK version, abstract SDK calls behind interface

- **Risk**: CGo memory leaks or crashes
  - **Likelihood**: Medium
  - **Impact**: High
  - **Mitigation**: Thorough testing, careful buffer management, ASAN testing

- **Risk**: USB bandwidth exhaustion with multiple streams
  - **Likelihood**: Low
  - **Impact**: Medium
  - **Mitigation**: Document bandwidth requirements, implement stream priority

## 10. Non-Functional Requirements

### Maintainability
- CGo bindings isolated in dedicated package
- Clear separation between SDK wrapper and business logic
- Comprehensive logging at DEBUG, INFO, WARN, ERROR levels

### Testability
- Mock camera interface for unit testing
- Integration tests with real hardware
- Benchmark tests for frame throughput

## 11. Future Considerations

### Phase 2 Features
- Multiple camera support (select by serial number)
- Hardware synchronization between cameras
- Depth-color alignment/registration

### Technical Debt
- Consider pure Go implementation if/when available
- Evaluate alternative SDKs (librealsense patterns)

---

**Document Status**: Draft

**Last Updated**: 2025-01-15

**Version**: 0.1
