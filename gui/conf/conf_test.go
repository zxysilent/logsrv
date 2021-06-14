package conf

import (
	"os"
	"testing"
)

func TestConf(t *testing.T) {
	defConfig = "./conf.toml"
	conf, err := initConf()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(conf)
	}
	os.RemoveAll("logs")
}
