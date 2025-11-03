# 安装指南

## 系统要求

- **操作系统**: Windows 10+, Linux, macOS
- **Go 版本**: 1.18 或更高
- **ffmpeg**: 必需

## 快速安装

### 1. 安装 ffmpeg

#### Windows
```powershell
# 使用 Chocolatey
choco install ffmpeg

# 验证
ffmpeg -version
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install ffmpeg

# 验证
ffmpeg -version
```

#### macOS
```bash
# 使用 Homebrew
brew install ffmpeg

# 验证
ffmpeg -version
```

### 2. 下载并运行

从 [Releases](https://github.com/your-repo/radikojp/releases) 页面下载最新版本，或者从源码编译：

```bash
# 克隆项目
git clone https://github.com/your-repo/radikojp.git
cd radikojp

# 安装依赖
go mod tidy

# 编译
go build -o radiko

# 运行
./radiko
```

## 详细说明

查看完整的安装指南，请参考项目根目录的文档。

---

**需要帮助？** 查看 [故障排除](TROUBLESHOOTING.md)
