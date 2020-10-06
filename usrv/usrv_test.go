package usrv

import (
	"os"
	"testing"
)

func TestTargz(t *testing.T) {
	t.Log(Targz("./../mock", "1"))
	os.RemoveAll("./logs")
}

func TestStat(t *testing.T) {
	_, err := os.Stat("./util.go")
	t.Log(err)
	_, err = os.Stat("./util")
	t.Log(err)
	_, err = os.Stat("./../conf")
	t.Log(err)
	os.RemoveAll("./logs")
}
