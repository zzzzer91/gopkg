package logx

import "testing"

func Test_Log(t *testing.T) {
	Debug("Debug")
	Debugf("Debug %s", "f")
	Info("Info")
	Info("Info %s", "f")
	Warn("Warn")
	Warnf("Warn %s", "f")
	Error("Error")
	Errorf("Error %s", "f")
	SetLevel(-1) // -1 represent debug level
	Debug("Debug")
	Debugf("Debug %s", "f")
	Error("Error")
	Errorf("Error %s", "f")
}
