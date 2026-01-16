Orbbec Femto Mega Network Streamer - Updated Project Spec
Purpose:
Go application that captures data from Orbbec Femto Mega and streams it over UDP network with runtime configuration control.
Stream Types:

Depth stream (ToF sensor data)
Color stream (RGB/4K video)
IR stream (optional)
IMU stream (6DoF, optional)

Per-Stream Configuration:

Enable/Disable - Turn stream on/off independently
Resolution - Available modes from Orbbec SDK (e.g., 640x576, 1021x1024 for depth; 1920x1080, 3840x2160 for color)
FPS - Frame rate (15/30/60fps depending on resolution)
Pixel Format - Depth: uint16, float; Color: RGB, RGBA, MJPEG, YUV

Network Transport:

Protocol: UDP
Port assignment: Separate port per stream type (configurable base port, e.g., 5000-5003)
Max UDP packet size: 65507 bytes

CRITICAL: Frame Chunking & Reassembly System
This is mission-critical for reliability - depth frames can be 2+ MB and require 30+ UDP packets
Chunking Strategy:

Packet header format: [Magic(4) | StreamType(1) | FrameID(4) | TotalChunks(2) | ChunkIndex(2) | Timestamp(8) | Width(2) | Height(2) | Format(1) | ChunkSize(4)] = 30 bytes header
Magic number for packet validation (e.g., 0xDEADBEEF)
Chunk size: 65,000 bytes payload (leaves room for header + network overhead)
Frame ID counter wraps at uint32 max

Reassembly Requirements:

Out-of-order packet handling - Chunks may arrive in any order
Missing packet detection - Timeout mechanism (e.g., 100ms per frame)
Frame buffer management - Ring buffer for N incomplete frames (suggest 3-5 frames)
Duplicate packet handling - Same chunk may arrive twice
Frame timeout/cleanup - Discard incomplete frames after deadline
Memory efficiency - Pre-allocate frame buffers, avoid allocations per-packet
Error correction - Option for FEC (Forward Error Correction) or simple parity chunks

Robustness Features:

CRC32 or checksum per chunk for corruption detection
Sequence number validation (detect dropped FrameIDs)
Bandwidth throttling to prevent network saturation
Automatic chunk size adjustment based on network conditions
Statistics tracking: packets sent/received, frames completed/dropped, reassembly time

Control Interface:

REST API on separate port (e.g., :8080) for runtime configuration
Endpoints:

GET /status - Current config + stream states + network stats
POST /streams/{type}/enable - Turn stream on/off
POST /streams/{type}/config - Update resolution/fps/format
GET /streams/{type} - Get current stream settings
GET /stats - Detailed chunking/reassembly statistics



Technical Requirements:

Use CGo bindings for Orbbec SDK (C/C++ API)
Goroutine per active stream (sender)
Separate goroutine for frame chunking/packetization
Efficient packet pooling to reduce GC pressure
Config file (YAML/JSON) for startup defaults
Graceful shutdown/cleanup
Error handling for camera disconnect/reconnect
Extensive logging for debugging chunking issues

Key Dependencies:

Orbbec SDK (C/C++ with CGo)
HTTP router (net/http or fiber/gin)
Config parser (viper or similar)
Consider: github.com/google/gopacket for packet utilities

Testing Requirements:

Simulated packet loss (5%, 10%, 20% drop rates)
Out-of-order delivery testing
High-load testing (multiple concurrent streams)
Network congestion simulation
Frame timeout validation

Bonus Features (Phase 2):

Multiple camera support (specify by serial number)
Frame compression (LZ4) per chunk
Multicast support
Adaptive bitrate based on network conditions
Optional TCP mode for guaranteed delivery (fallback)
WebSocket support for browser-based receivers

Documentation Needs:

Protocol specification document
Receiver implementation guide
Troubleshooting guide for packet loss scenarios