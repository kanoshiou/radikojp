package main

import (
	"fmt"
	"github.com/bluenviron/gohlslib/pkg/playlist"
	"io"
	"net/http"
	"os"
	"os/signal"
	"radikojp/hook"
	"radikojp/player"
	"syscall"
)

func main() {
	// 打印版本信息
	PrintVersion()
	
	url := "https://c-radiko.smartstream.ne.jp/QRR/_definst_/simul-stream.stream/playlist.m3u8?station_id=QRR&l=30&lsid=5e586af5ccb3b0b2498abfb19eaa8472&type=b"
	
	// 获取认证 token
	fmt.Println("Authenticating...")
	authToken := hook.Auth()
	fmt.Println("✓ Auth token obtained")

	// 获取播放列表
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Radiko-AuthToken", authToken)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// 解析播放列表
	byts, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	pl, err := playlist.Unmarshal(byts)
	if err != nil {
		panic(err)
	}

	streamUrl := ""

	switch pl := pl.(type) {
	case *playlist.Multivariant:
		fmt.Println("Multivariant playlist detected")
		if len(pl.Variants) > 0 {
			streamUrl = pl.Variants[0].URI
			fmt.Printf("Using stream: %s\n", streamUrl)
		}

	case *playlist.Media:
		fmt.Println("Media playlist detected")
		streamUrl = url
	}

	if streamUrl == "" {
		panic("No valid stream URL found")
	}

	// 创建并启动播放器
	fmt.Println("Starting ffmpeg player...")
	fmt.Println("Note: This requires ffmpeg to be installed and in PATH")
	fmt.Println()
	
	ffmpegPlayer := player.NewFFmpegPlayer(authToken)
	
	err = ffmpegPlayer.Play(streamUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to start player: %v", err))
	}

	fmt.Println()
	fmt.Println("🎵 Playing... Press Ctrl+C to stop")
	fmt.Println()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nStopping player...")
	ffmpegPlayer.Stop()
	fmt.Println("Stopped")
}
