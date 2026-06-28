package logger

import (
	"testing"
)

func TestInfoDoesNotPanic(t *testing.T) {
	Info("test %s", "info")
}

func TestSuccessDoesNotPanic(t *testing.T) {
	Success("test %s", "success")
}

func TestWarnDoesNotPanic(t *testing.T) {
	Warn("test %s", "warn")
}

func TestErrorDoesNotPanic(t *testing.T) {
	Error("test %s", "error")
}
