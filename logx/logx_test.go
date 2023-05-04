package logx

import "testing"

func Test_Log(t *testing.T) {
	Infoln("Info")
	Info("Info %s", "f")
	Warnln("Warn")
	Warnf("Warn %s", "f")
	Errorln("Error")
	Errorf("Error %s", "f")
}
