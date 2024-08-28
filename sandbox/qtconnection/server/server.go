package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Person struct {
	FirstName string
	LastName  string
	MyArray   []int
	MyMap     map[string]bool
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1200")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Connected: %s\n", conn.RemoteAddr().String())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer fmt.Printf("Disconnected: %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	for {
		// Read json object from client
		var cperson Person

		buf, _ := readAll(conn)
		err := json.Unmarshal(buf, &cperson)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print("From client: ")
		fmt.Println(cperson)

		// Write json object to client
		sperson := Person{
			FirstName: "Nancy",
			LastName:  "Star",
			MyArray:   []int{101, 23, 15, 66},
			MyMap:     map[string]bool{"101": true, "203": true},
		}

		data, err := json.Marshal(sperson)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Write(data)
	}
}

func readAll(conn net.Conn) ([]byte, error) {
	result := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if n < 512 {
			break
		}
	}

	return result.Bytes(), nil
}
