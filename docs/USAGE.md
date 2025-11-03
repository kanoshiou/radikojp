# 使用说明

## 基本使用

```bash
# 运行程序
./radiko

# 停止播放
按 Ctrl+C
```

## 自定义电台

编辑 `main.go` 修改电台 URL：

```go
url := "https://c-radiko.smartstream.ne.jp/TBS/_definst_/simul-stream.stream/playlist.m3u8?station_id=TBS&l=30&lsid=xxx&type=b"
```

## 配置选项

查看 `config.example.go` 了解可用的配置选项。

## 常见问题

### 没有声音？
- 检查系统音量
- 确认 ffmpeg 已安装
- 检查网络连接

### 播放卡顿？
- 检查网络速度
- 尝试更换网络

---

**更多信息** 查看 [故障排除](TROUBLESHOOTING.md)
