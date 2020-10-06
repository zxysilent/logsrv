package usrv

import (
	"encoding/json"
	"io/ioutil"
	"logsrv/conf"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"
	"github.com/zxysilent/logs"
)

var (
	pool   *sync.Pool //
	queue  chan Node  //缓冲队列
	task   *cron.Cron //
	global struct {
		Rows   int64 `json:"rows"` //
		Size   int64 `json:"size"` //
		SizeMB int64 `json:"size_mb"`
		SizeGB int64 `json:"size_gb"`
	}
	lock = &sync.Mutex{}
)

// Node Node
type Node struct {
	Msg  string
	Addr string
}

type applet map[string]*FileWriter

func (app applet) Get(key string) *FileWriter {
	fw, has := app[key]
	if has {
		return fw
	}
	fw, _ = NewWriter(conf.App.DefDir()+"/"+key+".log", key)
	app[key] = fw
	return fw
}
func (app applet) SyncFile() {
	for key := range app {
		app[key].SyncFile()
	}
}
func (app applet) Close() {
	for key := range app {
		app[key].Close()
	}
}

// Write 写入缓存
func Write(key, msg string) {
	lock.Lock()
	sb := app.Get(key)
	if sb != nil {
		ln, _ := sb.WriteString(msg + "\n")
		global.Size += int64(ln)
	}
	lock.Unlock()
}

// SyncFile 写入磁盘
func SyncFile() {
	logs.Info("sync file")
	lock.Lock()
	app.SyncFile()
	lock.Unlock()
}

// Rotate 备份
func Rotate() {
	lock.Lock()
	app.SyncFile()
	app.Close()
	now := time.Now()
	jbuf, _ := json.Marshal(app)
	os.Mkdir("json", 0666)
	ioutil.WriteFile(now.Format("json/2006-01-02150405.json"), jbuf, 0666)
	path := conf.App.RootDir + now.AddDate(0, 0, -1).Format("/2006-01-02150405")
	for key := range app {
		app[key].file.Close()
		delete(app, key)
	}
	logs.Info("rotate rename ", path, os.Rename(conf.App.DefDir(), path))
	lock.Unlock()
	logs.Info("rotate targz : ", path+conf.App.FileExt, Targz(path, path+conf.App.FileExt))
	logs.Info("rotate span : ", time.Now().Sub(now))
	logs.Info("delete old ", path, os.RemoveAll(path))

}

// Level 日志等级
// var Level = []string{"Emergency", "Alert", "Critical", "Error", "Warning", "Notice", "Informational", "Debug"}
var Level = []string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "debug"}
var app applet

func RunHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().Format("Server time => 2006-01-02 15:04:05")))
	})
	http.HandleFunc("/api/global", func(w http.ResponseWriter, r *http.Request) {
		global.SizeMB = global.Size / (1024 * 1024)
		global.SizeGB = global.SizeMB / 1024
		jbuf, _ := json.Marshal(global)
		w.Write(jbuf)
	})
	http.HandleFunc("/api/daily", func(w http.ResponseWriter, r *http.Request) {
		jbuf, _ := json.Marshal(app)
		w.Write(jbuf)
	})
	logs.Info("ListenHTTP on ", "[::]:80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		logs.Fatal("HTTP 监听失败!", err.Error())
	}
}
func RunUdp() {
	pool = &sync.Pool{New: func() interface{} { return make([]byte, 1536) }}
	queue = make(chan Node, 1e5)
	app = make(map[string]*FileWriter, 8)
	task = cron.New(cron.WithSeconds())
	logs.Debug(task.AddFunc("0 1 0 1/1 * ? ", Rotate))    //每天 00：01
	logs.Debug(task.AddFunc("0/30 * * * * ? ", SyncFile)) //每隔 30s
	task.Start()
	go solve()
	// 创建监听
	udp, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: conf.App.UdpPort})
	if err != nil {
		logs.Fatal("UDP 监听失败!", err.Error())
	}
	logs.Info("ListenUDP on ", udp.LocalAddr())
	go func() {
		for {
			buf := pool.Get().([]byte)
			ln, addr, err := udp.ReadFromUDP(buf)
			// 读取成功 放入队列
			if err == nil {
				node := Node{
					Msg:  string(buf[:ln]),
					Addr: addr.IP.String(),
				}
				queue <- node
			}
			pool.Put(buf)
			global.Rows++
			if global.Rows%conf.App.SplitSync == 0 {
				logs.Info("total-rows ", global.Rows, "queue ", len(queue))
				go SyncFile()
			}
		}
	}()
}

func solve() {
	for node := range queue { //等待关闭、并且数据为空
		e := strings.IndexByte(node.Msg, '>')
		if e < 1 {
			continue
		}
		pri, _ := strconv.Atoi(string(node.Msg[1:e]))
		key := node.Addr + "." + Level[pri%8]
		Write(key, node.Msg)
	}
}
