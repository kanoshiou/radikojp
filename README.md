# Radiko JP Player

一个用 Go 语言编写的 Radiko 日本网络电台播放器。

[![Release](https://img.shields.io/github/v/release/your-username/radikojp)](https://github.com/your-username/radikojp/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/your-username/radikojp)](https://go.dev/)
[![License](https://img.shields.io/github/license/your-username/radikojp)](LICENSE)

## ✨ 功能特性

- ✅ 支持 Radiko 认证
- ✅ 解析 HLS 播放列表
- ✅ 实时流式播放
- ✅ 跨平台支持（Windows/Linux/macOS）
- ✅ 使用 Go 原生音频输出库

## 📦 安装

### 方法 1: 下载预编译版本（推荐）

从 [Releases](https://github.com/your-username/radikojp/releases) 页面下载适合你系统的版本。

### 方法 2: 从源码编译

```bash
# 克隆项目
git clone https://github.com/your-username/radikojp.git
cd radikojp

# 安装依赖
go mod tidy

# 编译
go build -o radiko

# 运行
./radiko
```

## ⚠️ 重要提示

**需要安装 ffmpeg**：本程序使用 ffmpeg 进行 AAC 音频解码。

### 安装 ffmpeg

**Windows:**
```powershell
choco install ffmpeg
```

**Linux:**
```bash
sudo apt install ffmpeg  # Ubuntu/Debian
sudo yum install ffmpeg  # CentOS/RHEL
```

**macOS:**
```bash
brew install ffmpeg
```

验证安装：
```bash
ffmpeg -version
```

## 🚀 快速开始

```bash
# 运行程序
./radiko

# 停止播放
按 Ctrl+C
```

## 📖 文档

- [安装指南](docs/INSTALL.md)
- [使用说明](docs/USAGE.md)
- [故障排除](docs/TROUBLESHOOTING.md)

## 🏗️ 技术栈

- **HLS 处理**: [gohlslib](https://github.com/bluenviron/gohlslib)
- **音频输出**: [oto](https://github.com/hajimehoshi/oto)
- **音频解码**: ffmpeg

## 📋 系统要求

- Go 1.18+ （仅编译时需要）
- ffmpeg （运行时必需）
- 网络连接

## 🔧 开发

```bash
# 运行测试
go test ./...

# 格式化代码
go fmt ./...

# 检查代码
go vet ./...
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [gohlslib](https://github.com/bluenviron/gohlslib) - HLS 流处理
- [oto](https://github.com/hajimehoshi/oto) - 音频输出

---

**注意**: 本项目仅供学习和个人使用。请遵守 Radiko 的使用条款。
