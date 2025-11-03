.PHONY: all build run clean test fmt vet help

# 变量定义
BINARY_NAME=radiko
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 默认目标
all: build

# 编译
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

# 编译并运行
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# 清理
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@rm -f radiko-*
	@echo "Clean complete"

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
vet:
	@echo "Vetting code..."
	@go vet ./...

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# 跨平台编译
build-all:
	@echo "Building for all platforms..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-arm64.exe
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64
	@echo "Build complete for all platforms"

# 帮助信息
help:
	@echo "Radiko JP Player - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - 编译程序"
	@echo "  make run        - 编译并运行"
	@echo "  make clean      - 清理编译文件"
	@echo "  make test       - 运行测试"
	@echo "  make fmt        - 格式化代码"
	@echo "  make vet        - 检查代码"
	@echo "  make deps       - 下载依赖"
	@echo "  make build-all  - 跨平台编译"
	@echo "  make help       - 显示帮助信息"
