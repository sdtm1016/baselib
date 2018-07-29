package logger

import "testing"

func TestInfo(t *testing.T) {
	Info("hello","world")
	Debug("hello","world")
	Error("hello","world")
	Warn("hello","world")
}
