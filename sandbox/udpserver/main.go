package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listener, err := net.ListenPacket("udp", "127.0.0.1:1200")
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf [512]byte
	n, addr, err := listener.ReadFrom(buf[:])
	if err != nil {
		fmt.Println(err)
		listener.Close()
		return
	}
	fmt.Println(string(buf[:n]))

	listener.WriteTo([]byte("From server: "+time.Now().String()), addr)
	listener.Close()
}
