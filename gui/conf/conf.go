package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

type appconf struct {
	Title    string   `toml:"title"`
	Intro    string   `toml:"intro"`
	Mode     string   `toml:"mode"`
	Addr     string   `toml:"addr"`
	Srv      string   `toml:"srv"`
	VmIp     string   `toml:"vm_ip"`
	VmUser   string   `toml:"vm_user"`
	VmPasswd string   `toml:"vm_passwd"`
	VmDc     string   `toml:"vm_dc"`
	VmVms    []string `toml:"vm_vms"`
}

func (app *appconf) IsProd() bool {
	return app.Mode == "prod"
}
func (app *appconf) IsDev() bool {
	return app.Mode == "dev"
}

var (
	App       *appconf
	defConfig = "./conf/conf.toml"
)

func Init() {
	var err error
	App, err = initConf()
	if err != nil {
		logs.Fatal("config init error : ", err.Error())
	}
}

func initConf() (*appconf, error) {
	app := &appconf{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
