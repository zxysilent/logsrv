package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var host string
var port int

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "-h host")
	flag.IntVar(&port, "p", 514, "-p port")
}
func main() {
	flag.Parse()
	if len(os.Args) > 1 && os.Args[1] == "help" {
		flag.Usage()
		return
	}
	dstAddr := &net.UDPAddr{IP: net.ParseIP(host), Port: port}
	log.Println("send to :", dstAddr.String())
	time.Sleep(time.Second * 3)
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Println(err)
	}
	for {
		time.Sleep(time.Second * 1)
		// 发送数据
		conn.WriteToUDP([]byte("<132> test log warn test log warn test log warn test log warn"), dstAddr)
		ln, err := conn.WriteToUDP([]byte("<134> test log info test log info test log info test log info"), dstAddr)
		log.Println("send msg err:", err, ",length:", ln)
	}
}
