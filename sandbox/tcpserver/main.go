package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on port 8080
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listening: %s\n", listener.Addr().String())

	// Accept incoming connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Accepted: %s\n", conn.LocalAddr().String())

		// Handle the connect in a goroutine
		go handleConnect(conn)
	}
}

func handleConnect(conn net.Conn) {
	// Close when done
	defer conn.Close()

	for {
		// Read incoming data
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the incoming data
		fmt.Printf("Recieved: %s", buf)
	}
}
