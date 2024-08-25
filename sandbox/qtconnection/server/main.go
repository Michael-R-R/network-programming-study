package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type Person struct {
	FirstName string
	LastName  string
	MyArray   []int
	MyMap     map[int]bool
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

	var counter int

	for {
		var person Person

		buf, _ := readAll(conn)
		err := json.Unmarshal(buf, &person)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print("From client: ")
		fmt.Println(person)

		counter++
		conn.Write([]byte("Count: " + strconv.Itoa(counter)))
	}
}

func readAll(conn net.Conn) ([]byte, error) {
	resultBuf := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		resultBuf.Write(buf[0:n])

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

	return resultBuf.Bytes(), nil
}
