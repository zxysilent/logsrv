package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

type appconf struct {
	Title      string `toml:"title"`
	Intro      string `toml:"intro"`
	UdpPort    int    `toml:"upd_port"`
	HttpPort   int    `toml:"http_port"`
	HttpSrv    string `toml:"http_srv"`
	RootDir    string `toml:"root_dir"`
	FileExt    string `toml:"file_ext"`
	BufSize    int    `toml:"buf_size"` //KB
	SplitSync  int64  `toml:"split_sync"`
	SplitFile  int64  `toml:"split_file"`  //MB
	SaveUnknow bool   `toml:"save_unknow"` //MB
}

var (
	App       *appconf
	defConfig = "./conf/conf.toml"
)

func (app *appconf) DefDir() string {
	return app.RootDir + "/run"
}
func Init() {
	var err error
	App, err = initConf()
	if err != nil {
		logs.Fatal("config init error : ", err.Error())
	}
	App.BufSize *= 1024
	App.SplitFile *= 1024 * 1024
	logs.Debug("conf init")
}

func initConf() (*appconf, error) {
	app := &appconf{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
