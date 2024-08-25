package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1200")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write([]byte("From client: Hello, Server!"))
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	var buf [512]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	fmt.Println(string(buf[0:n]))

	conn.Close()
}
