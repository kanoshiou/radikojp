# Architecture

## Project Structure

```
radiko-tui/
├── .github/
│   ├── dependabot.yml            # Dependabot configuration
│   └── workflows/
│       ├── go-version-check.yml  # Go version compatibility check
│       └── release.yml           # GitHub Actions auto-release
├── api/
│   ├── auth.go                   # Radiko authentication module
│   └── client.go                 # Radiko API client
├── config/
│   └── config.go                 # Configuration management
├── docs/                         # Documentation directory
│   ├── ARCHITECTURE.md           # Architecture (this file)
│   ├── INSTALL.md                # Installation guide
│   ├── TROUBLESHOOTING.md        # Troubleshooting
│   └── USAGE.md                  # Usage guide
├── model/
│   ├── device.go                 # Device info and GPS generation
│   ├── program.go                # Program data models
│   ├── region.go                 # Region/Area definitions
│   └── station.go                # Station data models
├── player/
│   ├── ffmpeg_player.go          # FFmpeg-based audio player (with audio)
│   └── ffmpeg_player_noaudio.go  # Stub player (noaudio build)
├── server/
│   └── server.go                 # HTTP streaming server (StreamManager)
├── tui/
│   ├── tui.go                    # Terminal UI (with audio)
│   └── tui_noaudio.go            # Stub TUI (noaudio build)
├── main.go                       # Main program entry
├── config.example.go             # Configuration example
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── Makefile                      # Build script
├── README.md                     # Project description (English)
├── README.ja.md                  # Project description (Japanese)
├── README.zh.md                  # Project description (Chinese)
└── LICENSE                       # MIT License
```


## Technical Architecture

```
┌─────────────────────────────────────────────────────┐
│              Radiko JP Player                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────┐    ┌──────────────┐               │
│  │   TUI        │◄──►│   Config     │               │
│  │  (bubbletea) │    │   Manager    │               │
│  └──────────────┘    └──────────────┘               │
│         │                                           │
│         ▼                                           │
│  ┌──────────────┐    ┌──────────────┐               │
│  │ Auth Module  │───►│    Token     │               │
│  │  (auth1/2)   │    │   Manager    │               │
│  └──────────────┘    └──────────────┘               │
│         │                                           │
│         ▼                                           │
│  ┌──────────────┐    ┌──────────────┐               │
│  │  API Client  │───►│   Stream     │               │
│  │ (stations)   │    │   Selector   │               │
│  └──────────────┘    └──────────────┘               │
│         │                                           │
│         ▼                                           │
│  ┌──────────────┐    ┌──────────────┐               │
│  │   ffmpeg     │───►│  AAC→PCM     │               │
│  │  (external)  │    │   Decoder    │               │
│  └──────────────┘    └──────────────┘               │
│         │                                           │
│         ▼                                           │
│  ┌──────────────┐    ┌──────────────┐               │
│  │  oto Player  │───►│   Speaker    │               │
│  │    (Go)      │    │              │               │
│  └──────────────┘    └──────────────┘               │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## Core Modules

### 1. API Module (api/)

#### Authentication (api/auth.go)
Handles Radiko's two-step authentication:
- **auth1**: Obtains initial token, key offset, and length
- **auth2**: Validates token with partial key and GPS location
- Supports all 47 Japanese prefectures via GPS spoofing

#### API Client (api/client.go)
Communicates with Radiko services:
- `GetStations()`: Fetches station list for a region
- `GetStreamURLs()`: Gets streaming URLs for a station
- `GetCurrentProgram()`: Retrieves current program info
- `GetStationArea()`: Gets area ID for a station (auto-detection)

### 2. Player Module (player/ffmpeg_player.go)

FFmpeg-based audio player with:
- Real-time AAC to PCM decoding
- Software volume control
- Mute functionality
- Auto-reconnection on stream failure
- Reconnection status tracking

### 3. TUI Module (tui/tui.go)

Interactive terminal interface using bubbletea:
- Station list with scroll support
- Region selector (47 prefectures)
- Real-time volume display
- Current program display
- Keyboard navigation

### 4. Server Module (server/server.go)

HTTP streaming server for headless operation with advanced stream management:

#### Architecture
```
StreamManager
    └── StationStream (per station)
            ├── ffmpeg process
            ├── broadcast channel
            └── clients[] (multiple HTTP connections)
```

#### Features
- **Multi-client support**: Multiple clients can listen to the same station, sharing one ffmpeg instance
- **Smart ffmpeg reuse**: When a client disconnects, ffmpeg keeps running for a configurable grace period
- **Automatic reconnection**: If a client reconnects within the grace period, the existing stream is reused
- **Efficient broadcasting**: Data is read once from ffmpeg and broadcast to all connected clients

#### API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/play/{stationID}` | Stream audio from the specified station |
| `HEAD /api/play/{stationID}` | Get stream headers without starting playback |
| `GET /api/status` | Get JSON status of active streams |

#### Command Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-server` | false | Enable server mode |
| `-port` | 8080 | HTTP server port |
| `-grace` | 10 | Seconds to keep ffmpeg alive after last client disconnects |

Usage:
```bash
radiko-tui -server -port 8080 -grace 30
# Stream with: vlc http://localhost:8080/api/play/QRR
```

#### Server-Only Build (noaudio)

For headless Linux servers without audio support, build with the `noaudio` tag:
```bash
go build -tags noaudio -o radiko-server
```

This excludes the oto audio library and only supports server mode.

### 5. Configuration (config/config.go)

Persistent user preferences:
- Last played station
- Volume level
- Selected region
- Auto-saved on changes

### 6. Region/Device Models (model/)

- **region.go**: All 47 Japanese prefectures with IDs
- **device.go**: Random Android device generation for auth
- **program.go**: Program schedule data structures

## Data Flow

```
User Input (TUI)
    ↓
Region Selection → Auth Module (GPS spoofing)
    ↓
API Client → Station List
    ↓
User selects station
    ↓
API Client → Stream URLs
    ↓
FFmpeg Player ← HLS Stream (with auth token)
    ↓
AAC Decoding → PCM
    ↓
oto Library → Audio Output
```

## Dependencies

### Go Dependencies
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Terminal styling
- `github.com/charmbracelet/bubbles`: UI components
- `github.com/ebitengine/oto/v3`: Audio output

### External Dependencies
- `ffmpeg`: AAC audio decoding (required at runtime)

## Concurrency Model

- **Main goroutine**: TUI event loop
- **Audio pump goroutine**: Reads PCM from ffmpeg, writes to oto
- **Monitor goroutine**: Detects stream failures, triggers reconnect
- **ffmpeg process**: External process, communicates via stdout pipe

## Error Handling

- **Authentication failure**: Displays error in TUI, allows retry
- **Network error**: Auto-reconnects with new auth token
- **ffmpeg error**: Cleans up resources, shows error message
- **User interrupt**: Gracefully stops player and exits

## Configuration Storage

Config file location:
- **Windows**: `%APPDATA%\radiko-tui\config.json`
- **macOS**: `~/Library/Application Support/radiko-tui/config.json`
- **Linux**: `~/.config/radiko-tui/config.json`

## Performance

- **Memory**: ~20-30MB
- **CPU**: 5-15% (including ffmpeg)
- **Network**: ~128kbps (AAC stream)
- **Startup**: 2-3 seconds (auth + initial stream)
