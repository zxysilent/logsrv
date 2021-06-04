package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// dstAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 514}
	if len(os.Args) != 2 {
		log.Println("usage mock 127.0.0.1")
		return
	}
	dstAddr := &net.UDPAddr{IP: net.ParseIP(os.Args[1]), Port: 514}
	log.Println(dstAddr.String())
	time.Sleep(5 * time.Second)
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn.Write([]byte("hello"))
		fmt.Println(conn.RemoteAddr())
		time.Sleep(time.Millisecond * 5)
		// 发送数据
		conn.WriteToUDP([]byte("<132> test log warn test log warn test log warn test log warn"), dstAddr)
		ln, err := conn.WriteToUDP([]byte("<134> test log info test log info test log info test log info"), dstAddr)
		log.Println(ln, err)
	}
}
