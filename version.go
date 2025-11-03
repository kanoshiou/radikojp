package main

import "fmt"

var (
	// Version 程序版本（由构建时注入）
	Version = "dev"
	
	// BuildTime 构建时间（由构建时注入）
	BuildTime = "unknown"
)

// PrintVersion 打印版本信息
func PrintVersion() {
	fmt.Println("Radiko JP Player")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Println()
}
