# 架构说明

## 项目结构

```
radikojp/
├── .github/
│   └── workflows/
│       └── release.yml      # GitHub Actions 自动发布
├── docs/                    # 文档目录
│   ├── INSTALL.md          # 安装指南
│   ├── USAGE.md            # 使用说明
│   ├── TROUBLESHOOTING.md  # 故障排除
│   └── ARCHITECTURE.md     # 架构说明（本文件）
├── hook/
│   └── Auth.go             # Radiko 认证模块
├── model/
│   └── authtoken.go        # 数据模型
├── player/
│   └── ffmpeg_player.go    # FFmpeg 播放器实现
├── main.go                 # 主程序入口
├── version.go              # 版本信息
├── config.example.go       # 配置示例
├── go.mod                  # Go 模块定义
├── Makefile                # 构建脚本
├── README.md               # 项目说明
├── LICENSE                 # MIT 许可证
├── CHANGELOG.md            # 变更日志
└── CONTRIBUTING.md         # 贡献指南
```

## 技术架构

```
┌─────────────────────────────────────────────┐
│           Radiko JP Player                  │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────┐    ┌──────────┐             │
│  │  认证模块  │───→│ Token 管理│             │
│  └──────────┘    └──────────┘             │
│        │                                    │
│        ↓                                    │
│  ┌──────────┐    ┌──────────┐             │
│  │ HLS 解析  │───→│ 流选择    │             │
│  └──────────┘    └──────────┘             │
│        │                                    │
│        ↓                                    │
│  ┌──────────┐    ┌──────────┐             │
│  │  ffmpeg  │───→│ AAC 解码  │             │
│  │  (外部)   │    │  → PCM    │             │
│  └──────────┘    └──────────┘             │
│        │                                    │
│        ↓                                    │
│  ┌──────────┐    ┌──────────┐             │
│  │ oto 播放器│───→│  扬声器   │             │
│  │  (Go)    │    │          │             │
│  └──────────┘    └──────────┘             │
│                                             │
└─────────────────────────────────────────────┘
```

## 核心模块

### 1. 认证模块 (hook/Auth.go)

负责 Radiko 的认证流程：
- auth1: 获取初始 token
- auth2: 验证并激活 token

### 2. 播放器模块 (player/ffmpeg_player.go)

使用 ffmpeg 解码 AAC 音频：
- 创建 ffmpeg 子进程
- 通过管道传输 PCM 数据
- 使用 oto 库输出音频

### 3. 主程序 (main.go)

协调各个模块：
- 获取认证
- 解析播放列表
- 启动播放器
- 处理信号

## 数据流

```
Radiko API
    ↓ (HTTP + Auth Token)
HLS Playlist (m3u8)
    ↓ (解析)
Stream URL
    ↓ (HTTP + Auth Token)
AAC Audio Stream
    ↓ (ffmpeg 解码)
PCM Audio Data
    ↓ (oto 播放)
扬声器输出
```

## 依赖关系

### Go 依赖
- `gohlslib`: HLS 播放列表解析
- `oto/v2`: 跨平台音频输出
- `beep`: 音频处理工具

### 外部依赖
- `ffmpeg`: AAC 音频解码

## 并发模型

- 主 goroutine: 处理用户输入和信号
- 播放 goroutine: 管理 ffmpeg 进程和音频输出
- ffmpeg 进程: 独立进程，通过管道通信

## 错误处理

- 认证失败: 打印错误并退出
- 网络错误: 打印错误并退出
- ffmpeg 错误: 打印错误并清理资源
- 用户中断: 优雅关闭所有资源

## 性能考虑

- **内存**: ~20-30MB
- **CPU**: 5-15% (包括 ffmpeg)
- **网络**: ~128kbps
- **延迟**: 2-3秒启动延迟

## 未来改进

1. **纯 Go AAC 解码器**: 移除 ffmpeg 依赖
2. **缓冲优化**: 改进缓冲策略
3. **错误恢复**: 自动重连机制
4. **多电台支持**: 命令行参数选择电台
5. **GUI**: 图形界面

## 参考资料

- [Radiko API](https://radiko.jp)
- [HLS 协议](https://datatracker.ietf.org/doc/html/rfc8216)
- [AAC 格式](https://wiki.multimedia.cx/index.php/ADTS)
- [oto 文档](https://pkg.go.dev/github.com/hajimehoshi/oto/v2)
