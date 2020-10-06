package main

import (
	"flag"
	"logsrv/conf"
	"logsrv/usrv"
	"os"
	"os/signal"
	"syscall"

	"github.com/zxysilent/logs"
)

// logs.SetLevel(logs.WARN)
// logs.SetLevel(logs.DEBUG)
// logs.SetCallInfo(true)
// logs.SetConsole(true)
func main() {
	flag.Parse()
	logs.SetLevel(logs.DEBUG)
	logs.SetCallInfo(true)
	logs.SetConsole(true)
	logs.Info("app initializing")
	conf.Init()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	go usrv.RunHttp()
	go usrv.RunUdp()
	<-quit
	logs.Info("app quitted")
	logs.Flush()
}
