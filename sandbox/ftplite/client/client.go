package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const ( // Message types
	LS    = "ls"
	CD    = "cd"
	PWD   = "pwd"
	EXIT  = "exit"
	ERROR = "error"
)

type Packet struct {
	Keys   []string
	Values []string
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1200")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		vals := strings.Split(cmd, " ")
		if len(vals) < 1 {
			continue
		}

		cmd = strings.TrimSpace(vals[0])
		cmd = strings.ToLower(cmd)

		switch cmd {
		case LS:
			{
				handleLS(cmd, conn)
			}
		case EXIT:
			{
				handleExit(cmd, conn)
				return
			}
		default:
			{
				fmt.Printf("invalid command: %s\n", cmd)
				continue
			}
		}
	}
}

func handleLS(cmd string, conn net.Conn) {
	// Write to server
	cpacket := Packet{Keys: []string{cmd}, Values: []string{""}}
	data, _ := json.Marshal(cpacket)
	conn.Write(data)

	// Read server response
	data, err := readAll(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	var spacket Packet
	json.Unmarshal(data, &spacket)

	if spacket.Keys[0] == ERROR {
		fmt.Println(spacket.Values[0])
		return
	}

	dlist := strings.Split(spacket.Values[0], ",")

	fmt.Println("---------")
	for _, d := range dlist {
		fmt.Println(d)
	}
	fmt.Println("---------")
}

func handleExit(cmd string, conn net.Conn) {
	packet := Packet{Keys: []string{cmd}, Values: []string{""}}
	data, _ := json.Marshal(packet)

	conn.Write(data)

	fmt.Printf("Closing connection: %s\n", conn.RemoteAddr().String())
}

func readAll(conn net.Conn) ([]byte, error) {
	result := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[:])
		result.Write(buf[:n])

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
