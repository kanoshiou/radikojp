# Radiko TUI

[English](README.md) | **[日本語](README.ja.md)** | [中文](README.zh.md)

Go言語で書かれたRadiko日本インターネットラジオのターミナルUI（TUI）プレーヤーです。

[![Release](https://img.shields.io/github/v/release/kanoshiou/radiko-tui)](https://github.com/kanoshiou/radiko-tui/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kanoshiou/radiko-tui)](https://go.dev/)
[![License](https://img.shields.io/github/license/kanoshiou/radiko-tui)](LICENSE)

## ✨ 機能

- 🎵 Radikoラジオ局のライブストリーミング
- 🗾 日本全国47都道府県対応
- 🖥️ インタラクティブなターミナルUI（TUI）
- 🌐 HTTPストリーミングのサーバーモード
- 🔊 ミュート機能付き音量調整
- ⏺️ AACファイルへのストリーム録音
- 🔄 ストリーム障害時の自動再接続
- 💾 前回の放送局と設定を記憶
- 🌏 クロスプラットフォーム（Windows/Linux/macOS）

## 📸 スクリーンショット

```
📻 Radiko  🔊 80%
  ◀ 埼玉 千葉 [東京] 神奈川 新潟 ▶ [13/47]
──────────────────────────────────────────────
  TBSラジオ TBS
 ▶ 文化放送 QRR 
  ニッポン放送 LFR
  ラジオNIKKEI第1 RN1
  ラジオNIKKEI第2 RN2
  ↓ さらに表示

──────────────────────────────────────────────
▶ 文化放送 QRR  ♪ 大竹まことゴールデンラジオ  ⏺ 録音中 02:15
↑↓ 選択  Enter 再生  ←→ 地域切替  +- 音量  m ミュート  s 停止  r 再接続  Esc 終了
```

## 📦 インストール

### ビルド済みバイナリのダウンロード（推奨）

[Releases](https://github.com/kanoshiou/radiko-tui/releases) からダウンロードしてください。

### ソースからビルド

```bash
git clone https://github.com/kanoshiou/radiko-tui.git
cd radiko-tui
go mod tidy
go build -o radiko
```

## ⚠️ 必要条件

音声デコードと録音には **ffmpeg が必要** です。

```bash
# Windows (Chocolatey)
choco install ffmpeg

# Linux (Ubuntu/Debian)
sudo apt install ffmpeg

# macOS (Homebrew)
brew install ffmpeg
```

## 🚀 使用方法

### TUIモード（デフォルト）

```bash
./radiko-tui
```

### サーバーモード

HTTPストリーミングサーバーとして実行：

```bash
./radiko-tui -server -port 8080
```

VLCまたは任意のオーディオプレーヤーでストリーミング：

```bash
vlc http://localhost:8080/api/play/QRR
```

### 操作方法

| キー | 操作 |
|-----|--------|
| ↑/↓ または k/j | 放送局を選択 |
| ←/→ または h/l | 地域を切り替え |
| Enter/Space | 放送局を再生 |
| +/- | 音量調整 |
| 0-9 | 音量レベルを設定 |
| m | ミュート切り替え |
| s | 録音開始/停止 |
| r | 再接続 |
| Esc | 終了 |

### 録音機能

`s` キーを押すと、現在のストリームの録音を開始/停止できます。録音ファイルはダウンロードフォルダに `radiko_放送局名_YYYYMMDD_HHMMSS.aac` の形式で保存されます。

再生中の放送局と異なる放送局を録音している場合、放送局名が括弧で表示されます：`⏺ 録音中[放送局名] MM:SS`

## 📖 ドキュメント

- [インストールガイド](docs/INSTALL.md)
- [使用ガイド](docs/USAGE.md)
- [トラブルシューティング](docs/TROUBLESHOOTING.md)
- [アーキテクチャ](docs/ARCHITECTURE.md)

## 🏗️ 技術スタック

- **TUI**: [bubbletea](https://github.com/charmbracelet/bubbletea)
- **オーディオ**: [oto](https://github.com/ebitengine/oto) + ffmpeg
- **スタイリング**: [lipgloss](https://github.com/charmbracelet/lipgloss)

## 🙏 謝辞

インスピレーションと参考のために [rajiko](https://github.com/jackyzy823/rajiko) に感謝します。

## 📋 システム要件

- ffmpeg（実行時）
- Go 1.18+（ビルド時のみ）
- UTF-8対応ターミナル

## 🤝 貢献

IssueおよびPull Requestを歓迎します！

## 📄 ライセンス

MITライセンス - [LICENSE](LICENSE) を参照

---

**注意**: このプロジェクトは学習および個人使用のみを目的としています。Radikoの利用規約を遵守してください。
