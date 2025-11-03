# 故障排除

## 常见问题

### ffmpeg not found
确保 ffmpeg 已安装并在 PATH 中：
```bash
ffmpeg -version
```

### 没有声音
1. 检查系统音量
2. 确认音频设备正常
3. 查看程序输出的错误信息

### 编译错误
```bash
# 清理并重新下载依赖
go clean -modcache
go mod tidy
go build
```

## 技术说明

本程序使用 ffmpeg 解码 AAC 音频流。AAC 是压缩格式，必须解码成 PCM 才能播放。

## 获取帮助

- 查看 [GitHub Issues](https://github.com/your-repo/radikojp/issues)
- 提交新的 Issue
