package main

import "fmt"

var (
	// Version is the program version (injected at build time)
	Version = "dev"
	
	// BuildTime is the build time (injected at build time)
	BuildTime = "unknown"
)

// PrintVersion prints version information
func PrintVersion() {
	fmt.Println("Radiko JP Player")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Println()
}
