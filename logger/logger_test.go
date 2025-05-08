package logger

import "testing"

func TestLogger(t *testing.T) {
	Initialize(ConsoleEncode, DebugLevel)

	Infof("test %s", "test")
	Debug("test")
	Errorf("test %s", "test")
}
