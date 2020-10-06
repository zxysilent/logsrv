package usrv

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"logsrv/conf"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zxysilent/logs"
)

// FileWriter doc
// 缓存写文件
// 每个ip每种日志一个FileWriter
type FileWriter struct {
	file    *os.File
	bwr     *bufio.Writer
	Rows    int64     `json:"rows"` // 累计行数
	Size    int64     `json:"size"` // 累计大小
	fsize   int64     // 文件大小
	frows   int64     // 文件行数
	created time.Time // 文件创建日期
	fpath   string    // 文件目录 完整路径 fpath=fname+fsuffix
	fname   string    // 文件名
	key     string    //
	fsuffix string    // 文件后缀名 默认 .log
}

// @param name 根目录+name
func NewWriter(name string, key string) (*FileWriter, error) {
	fi, err := os.Stat(name)
	fw := new(FileWriter)
	fw.created = time.Now()
	fw.key = key
	// logs/app.log
	fw.fsuffix = filepath.Ext(name)                 // .log
	fw.fname = strings.TrimSuffix(name, fw.fsuffix) // logs/app
	if fw.fsuffix == "" {
		fw.fsuffix = ".log"
	}
	fw.fpath = fw.fname + fw.fsuffix
	os.MkdirAll(filepath.Dir(name), 0666)
	if err == nil && !fi.IsDir() { //文件存在、不是目录
		fw.created = fi.ModTime()
		fw.Size = fi.Size()
		fw.fsize = fw.Size
		fw.Rows = ReadLines(name)
		fw.frows = fw.Rows
	}
	file, err := os.OpenFile(fw.fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logs.Error("NewWriter open file : ", err.Error())
		return nil, err
	}
	fw.file = file
	fw.bwr = bufio.NewWriterSize(fw.file, conf.App.BufSize) //KB
	return fw, nil
}
func (fw *FileWriter) WriteString(msg string) (int, error) {
	fw.Size += int64(len(msg)) //累计
	fw.fsize += int64(len(msg))
	fw.Rows += 1 //累计
	fw.frows += 1
	if fw.fsize >= conf.App.SplitFile { //MB
		fw.Rotate()
	}
	return fw.bwr.WriteString(msg)
}
func (fw *FileWriter) State() (int64, int64) {
	return fw.Size, fw.Rows
}

// SyncFile 写入磁盘
func (fw *FileWriter) SyncFile() {
	fw.bwr.Flush()
	fw.file.Sync()
}

// Close 关闭文件
func (fw *FileWriter) Close() {
	fw.SyncFile()
	fw.file.Close()
}

// rotate
func (fw *FileWriter) Rotate() error {
	// 存档
	part := strconv.FormatInt(fw.Size/conf.App.SplitFile, 10)
	if fw.file != nil {
		fw.bwr.Flush()
		fw.file.Sync()
		fw.file.Close()
		// 保存
		fbak := filepath.Join(fw.fname + ".part." + part + fw.fsuffix)
		os.Rename(fw.fpath, fbak)
	}
	// 新建
	file, err := os.OpenFile(fw.fpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logs.Error("rotate open file : ", err.Error())
		return err
	}
	fw.frows = 0
	fw.fsize = 0
	fw.file = file
	fw.created = time.Now()
	fw.bwr = bufio.NewWriterSize(fw.file, conf.App.BufSize) //KB
	logs.Info("rotate file-writer file")
	return nil
}
func ReadLines(name string) int64 {
	file, err := os.Open(name)
	if err != nil {
		return 0
	}
	defer file.Close()
	fd := bufio.NewReader(file)
	rows := int64(0)
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		rows++
	}
	return rows
}

type State struct {
	Date      time.Time        `json:"date"`
	TotalRows int64            `json:"total_rows"` //累计
	TotalSize int64            `json:"total_size"` //累计
	Rows      int64            `json:"rows"`       //当天
	Size      int64            `json:"size"`       //当天
	Items     map[string]*Item `json:"items"`
}
type Item struct {
	Name string `json:"name"`
	Rows int64  `json:"rows"`
	Size int64  `json:"size"`
}

// Targz 压缩文件 .tar.gz
func Targz(src string, dst string) error {
	if !strings.HasSuffix(dst, ".tar.gz") {
		dst += ".tar.gz"
	}
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			file.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
