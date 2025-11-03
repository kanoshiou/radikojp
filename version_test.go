package main

import (
	"testing"
)

func TestPrintVersion(t *testing.T) {
	// 测试版本打印不会panic
	PrintVersion()
}

func TestVersionVariables(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	
	if BuildTime == "" {
		t.Error("BuildTime should not be empty")
	}
}
