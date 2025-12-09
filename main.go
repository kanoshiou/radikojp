package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"radikojp/api"
	"radikojp/hook"
	"radikojp/player"
	"syscall"
	"time"

	"github.com/bluenviron/gohlslib/pkg/playlist"
	"github.com/eiannone/keyboard"
)

func main() {
	// 解析命令行参数
	volumePercent := flag.Int("volume", 80, "Initial volume (0-100)")
	flag.Parse()

	// 转换为 0.0-1.0 范围
	initialVolume := float64(*volumePercent) / 100.0
	if initialVolume < 0 {
		initialVolume = 0
	} else if initialVolume > 1 {
		initialVolume = 1
	}

	// 获取认证 token
	fmt.Println("Authenticating...")
	authToken := hook.Auth()
	fmt.Println("✓ Auth token obtained")

	// 获取电台列表
	fmt.Println("Fetching station list...")
	stations, err := api.GetStations()
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch stations: %v", err))
	}
	fmt.Printf("✓ Found %d stations\n", len(stations))

	if len(stations) == 0 {
		panic("No stations available")
	}

	// 默认选择 QRR 或第一个
	currentStationIdx := 0
	for i, s := range stations {
		if s.ID == "QRR" {
			currentStationIdx = i
			break
		}
	}

	// 创建播放器
	fmt.Println("Starting ffmpeg player...")
	fmt.Println("Note: This requires ffmpeg to be installed and in PATH")
	fmt.Printf("Initial volume: %d%%\n", *volumePercent)
	fmt.Println()

	ffmpegPlayer := player.NewFFmpegPlayer(authToken, initialVolume)

	// 设置重连回调函数
	ffmpegPlayer.SetReconnectCallback(func() string {
		return hook.Auth()
	})

	// 播放指定电台的函数
	// 播放指定电台的函数
	playStation := func(idx int) {
		station := stations[idx]
		fmt.Printf("\n\n📺 Switching to: %s (%s)\n", station.Name, station.ID)

		// 获取播放列表 URL 列表
		playlistURLs, err := api.GetStreamURLs(station.ID)
		if err != nil {
			fmt.Printf("❌ Failed to get stream URLs: %v\n", err)
			return
		}

		var finalStreamUrl string
		lsid := "5e586af5ccb3b0b2498abfb19eaa8472"
		
		// 使用最后一个 URL（与之前的行为一致）
		if len(playlistURLs) > 0 {
			lastUrl := playlistURLs[len(playlistURLs)-1]
			finalStreamUrl = fmt.Sprintf("%s?station_id=%s&l=30&lsid=%s&type=b", lastUrl, station.ID, lsid)
		}

		if finalStreamUrl == "" {
			fmt.Printf("❌ No stream URL available for station %s\n", station.Name)
			return
		}

		// 停止当前播放（如果有）
		ffmpegPlayer.Stop()
		
		println(finalStreamUrl)
		// 启动新播放
		err = ffmpegPlayer.Play(finalStreamUrl)
		if err != nil {
			fmt.Printf("❌ Failed to start player: %v\n", err)
			return
		}

		// 等待播放器启动
		time.Sleep(500 * time.Millisecond)
		fmt.Println("🎵 Playing...")
		printControls()
		printVolumeStatus(ffmpegPlayer, station.Name)
	}

	// 初始播放
	playStation(currentStationIdx)

	// 初始化键盘监听
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Warning: Could not open keyboard: %v\n", err)
		fmt.Println("Controls disabled. Press Ctrl+C to stop")

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
	} else {
		defer keyboard.Close()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// 键盘处理循环
		go func() {
			lastUpdate := time.Now()
			updateInterval := 50 * time.Millisecond

			for {
				char, key, err := keyboard.GetKey()
				if err != nil {
					return
				}

				needsUpdate := false
				needsSwitch := false

				switch key {
				case keyboard.KeyArrowUp:
					ffmpegPlayer.IncreaseVolume(0.05)
					needsUpdate = true
				case keyboard.KeyArrowDown:
					ffmpegPlayer.DecreaseVolume(0.05)
					needsUpdate = true
				case keyboard.KeyArrowRight: // 下一个电台
					currentStationIdx = (currentStationIdx + 1) % len(stations)
					needsSwitch = true
				case keyboard.KeyArrowLeft: // 上一个电台
					currentStationIdx = (currentStationIdx - 1 + len(stations)) % len(stations)
					needsSwitch = true
				}

				switch char {
				case '+', '=':
					ffmpegPlayer.IncreaseVolume(0.05)
					needsUpdate = true
				case '-', '_':
					ffmpegPlayer.DecreaseVolume(0.05)
					needsUpdate = true
				case 'e', 'E':
					ffmpegPlayer.IncreaseVolume(0.05)
					needsUpdate = true
				case 'q', 'Q':
					ffmpegPlayer.DecreaseVolume(0.05)
					needsUpdate = true
				case 'm', 'M':
					ffmpegPlayer.ToggleMute()
					needsUpdate = true
				case 'n', 'N': // Next station
					currentStationIdx = (currentStationIdx + 1) % len(stations)
					needsSwitch = true
				case 'p', 'P': // Previous station
					currentStationIdx = (currentStationIdx - 1 + len(stations)) % len(stations)
					needsSwitch = true
				case 'r', 'R':
					fmt.Println("\n🔄 Reconnecting...")
					go func() {
						err := ffmpegPlayer.Reconnect()
						if err != nil {
							fmt.Printf("\n❌ Reconnect failed: %v\n", err)
						}
					}()
					continue
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					volume := float64(char-'0') / 10.0
					ffmpegPlayer.SetVolume(volume)
					needsUpdate = true
				}

				if needsSwitch {
					playStation(currentStationIdx)
				} else if needsUpdate && time.Since(lastUpdate) > updateInterval {
					printVolumeStatus(ffmpegPlayer, stations[currentStationIdx].Name)
					lastUpdate = time.Now()
				}
			}
		}()

		<-sigChan
	}

	fmt.Println("\nStopping player...")
	ffmpegPlayer.Stop()
	fmt.Println("Stopped")
}

func resolveStreamURL(playlistURL, authToken string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", playlistURL, nil)
	req.Header.Set("X-Radiko-AuthToken", authToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	byts, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	pl, err := playlist.Unmarshal(byts)
	if err != nil {
		return "", err
	}

	switch pl := pl.(type) {
	case *playlist.Multivariant:
		if len(pl.Variants) > 0 {
			return pl.Variants[0].URI, nil
		}
		return "", fmt.Errorf("no variants in playlist")
	case *playlist.Media:
		return playlistURL, nil
	default:
		return "", fmt.Errorf("unknown playlist type")
	}
}

func printControls() {
	fmt.Println("Controls:")
	fmt.Println("  ↑ / + / e     Increase volume")
	fmt.Println("  ↓ / - / q     Decrease volume")
	fmt.Println("  → / n         Next Station")
	fmt.Println("  ← / p         Previous Station")
	fmt.Println("  m             Mute/Unmute")
	fmt.Println("  r             Reconnect/Replay")
	fmt.Println("  0-9           Set volume to 0%-90%")
	fmt.Println("  Ctrl+C        Stop and exit")
	fmt.Println()
}

func printVolumeStatus(p *player.FFmpegPlayer, stationName string) {
	volume := int(p.GetVolume() * 100)
	muted := p.IsMuted()

	status := fmt.Sprintf("[%s] Vol: %3d%%", stationName, volume)
	if muted {
		status += " [MUTED]"
	} else {
		status += "        "
	}

	barLength := 20
	filledLength := int(float64(barLength) * p.GetVolume())
	bar := ""
	for i := 0; i < barLength; i++ {
		if i < filledLength && !muted {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	output := fmt.Sprintf("%s [%s]", status, bar)
	// 使用 \r 和空格清除行，防止残留
	fmt.Printf("\r%-80s", output)
}
