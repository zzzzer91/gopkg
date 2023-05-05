package logx

import "testing"

func Test_Log(t *testing.T) {
	Debugln("Debug")
	Debugf("Debug %s", "f")
	Infoln("Info")
	Info("Info %s", "f")
	Warnln("Warn")
	Warnf("Warn %s", "f")
	Errorln("Error")
	Errorf("Error %s", "f")
	SetLevel(-1) // -1 represent debug level
	Debugln("Debug")
	Debugf("Debug %s", "f")
	Errorln("Error")
	Errorf("Error %s", "f")
}
