package player

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

// FFmpegPlayer 使用 ffmpeg 解码的播放器
type FFmpegPlayer struct {
	authToken   string
	mu          sync.Mutex
	playing     bool
	ctx         context.Context
	cancel      context.CancelFunc
	cmd         *exec.Cmd
	otoContext  *oto.Context
	otoPlayer   oto.Player
}

// NewFFmpegPlayer 创建 ffmpeg 播放器
func NewFFmpegPlayer(authToken string) *FFmpegPlayer {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &FFmpegPlayer{
		authToken: authToken,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Play 开始播放
func (p *FFmpegPlayer) Play(streamURL string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		return fmt.Errorf("already playing")
	}

	// 检查 ffmpeg 是否可用
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found in PATH. Please install ffmpeg: %w", err)
	}

	// 初始化音频输出
	err = p.initAudio(48000, 2)
	if err != nil {
		return fmt.Errorf("failed to init audio: %w", err)
	}

	// 创建 ffmpeg 命令
	// 使用 ffmpeg 解码 HLS 流并输出 PCM
	p.cmd = exec.CommandContext(p.ctx, "ffmpeg",
		"-headers", fmt.Sprintf("X-Radiko-AuthToken: %s", p.authToken),
		"-i", streamURL,
		"-f", "s16le",      // 16-bit PCM
		"-ar", "48000",     // 48kHz
		"-ac", "2",         // 立体声
		"-loglevel", "error", // 只显示错误
		"pipe:1",           // 输出到 stdout
	)

	// 获取 stdout
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// 启动 ffmpeg
	fmt.Println("Starting ffmpeg decoder...")
	err = p.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	p.playing = true

	// 启动数据传输
	go p.pumpAudio(stdout)

	fmt.Println("Playback started")
	return nil
}

// initAudio 初始化音频输出
func (p *FFmpegPlayer) initAudio(sampleRate, channelCount int) error {
	fmt.Printf("Initializing audio: %d Hz, %d channels\n", sampleRate, channelCount)

	var ready chan struct{}
	var err error
	p.otoContext, ready, err = oto.NewContext(sampleRate, channelCount, oto.FormatSignedInt16LE)
	if err != nil {
		return fmt.Errorf("failed to create oto context: %w", err)
	}

	<-ready
	fmt.Println("Audio context ready")
	return nil
}

// pumpAudio 传输音频数据
func (p *FFmpegPlayer) pumpAudio(reader io.Reader) {
	// 创建播放器
	p.otoPlayer = p.otoContext.NewPlayer(reader)
	p.otoPlayer.Play()

	fmt.Println("Audio pump started")

	// 等待上下文取消
	<-p.ctx.Done()
	
	fmt.Println("Audio pump stopped")
}

// Stop 停止播放
func (p *FFmpegPlayer) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return
	}

	fmt.Println("Stopping ffmpeg player...")
	
	p.cancel()
	
	if p.otoPlayer != nil {
		p.otoPlayer.Close()
	}
	
	if p.otoContext != nil {
		p.otoContext.Suspend()
	}
	
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
	}

	p.playing = false
	fmt.Println("FFmpeg player stopped")
}

// IsPlaying 是否正在播放
func (p *FFmpegPlayer) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing
}
