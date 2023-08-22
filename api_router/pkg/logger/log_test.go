package logger

import (
	"testing"
)

func TestLog(t *testing.T) {

	//_, filePath, _, _ := runtime.Caller(0)
	//join := filepath.Join(filePath, "..", "..", "..", "logs")
	//print(join)

	Log.Info("222")
}
