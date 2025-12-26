//go:build noaudio

package player

import (
	"fmt"
	"time"
)

// FFmpegPlayer is a stub player for server-only builds without audio support
type FFmpegPlayer struct {
	authToken string
	streamURL string
	volume    float64
}

// ReconnectStatus represents the reconnection state
type ReconnectStatus int

const (
	ReconnectNone ReconnectStatus = iota
	ReconnectStarted
	ReconnectAuth
	ReconnectPlaying
	ReconnectSuccess
	ReconnectFailed
)

// NewFFmpegPlayer creates a new ffmpeg player stub
func NewFFmpegPlayer(authToken string, initialVolume float64) *FFmpegPlayer {
	return &FFmpegPlayer{
		authToken: authToken,
		volume:    initialVolume,
	}
}

// SetReconnectCallback is a no-op in server-only mode
func (p *FFmpegPlayer) SetReconnectCallback(callback func() string) {}

// UpdateAuthToken updates the authentication token
func (p *FFmpegPlayer) UpdateAuthToken(token string) {
	p.authToken = token
}

// GetReconnectStatus returns the current reconnection status
func (p *FFmpegPlayer) GetReconnectStatus() ReconnectStatus {
	return ReconnectNone
}

// GetLastError returns the last error message
func (p *FFmpegPlayer) GetLastError() string {
	return ""
}

// ClearReconnectStatus clears the reconnection status
func (p *FFmpegPlayer) ClearReconnectStatus() {}

// Play returns an error since audio is not supported
func (p *FFmpegPlayer) Play(streamURL string) error {
	return fmt.Errorf("音声再生はサポートされていません (noaudio build)")
}

// Stop is a no-op in server-only mode
func (p *FFmpegPlayer) Stop() {}

// IsPlaying always returns false in server-only mode
func (p *FFmpegPlayer) IsPlaying() bool {
	return false
}

// SetVolume is a no-op in server-only mode
func (p *FFmpegPlayer) SetVolume(volume float64) {
	p.volume = volume
}

// GetVolume returns the current volume
func (p *FFmpegPlayer) GetVolume() float64 {
	return p.volume
}

// IncreaseVolume is a no-op in server-only mode
func (p *FFmpegPlayer) IncreaseVolume(delta float64) {}

// DecreaseVolume is a no-op in server-only mode
func (p *FFmpegPlayer) DecreaseVolume(delta float64) {}

// ToggleMute is a no-op in server-only mode
func (p *FFmpegPlayer) ToggleMute() {}

// IsMuted always returns false in server-only mode
func (p *FFmpegPlayer) IsMuted() bool {
	return false
}

// Reconnect is not supported in server-only mode
func (p *FFmpegPlayer) Reconnect() error {
	return fmt.Errorf("再接続はサポートされていません (noaudio build)")
}

// StartRecording is not supported in server-only mode
func (p *FFmpegPlayer) StartRecording(stationName string) error {
	return fmt.Errorf("録音はサポートされていません (noaudio build)")
}

// StopRecording is not supported in server-only mode
func (p *FFmpegPlayer) StopRecording() (string, error) {
	return "", fmt.Errorf("録音はサポートされていません (noaudio build)")
}

// IsRecording always returns false in server-only mode
func (p *FFmpegPlayer) IsRecording() bool {
	return false
}

// GetRecordingInfo returns empty values in server-only mode
func (p *FFmpegPlayer) GetRecordingInfo() (filePath string, duration time.Duration, stationName string) {
	return "", 0, ""
}

// ToggleRecording is not supported in server-only mode
func (p *FFmpegPlayer) ToggleRecording(stationName string) (started bool, filePath string, err error) {
	return false, "", fmt.Errorf("録音はサポートされていません (noaudio build)")
}
