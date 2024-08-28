package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection
	defer conn.Close()

	for {
		// Ask for input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message: ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Send some data to the server
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
