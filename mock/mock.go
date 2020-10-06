package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// dstAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 514}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("192.168.149.131"), Port: 514}
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
