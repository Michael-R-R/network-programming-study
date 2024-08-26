package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type User struct {
	WorkingDir string
}

type DirList struct {
	Directories []string
}

type ErrorMsg struct {
	Error string
}

var connections map[string]net.Conn

func init() {
	connections = make(map[string]net.Conn)
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer delete(connections, conn.RemoteAddr().String())
	defer conn.Close()

	user := User{WorkingDir: "."}

	connections[conn.RemoteAddr().String()] = conn

	var buf [512]byte

	for {
		// Read msg type
		n, err := conn.Read(buf[:])
		if err != nil {
			handleError(err, conn)
			continue
		}

		msg := strings.TrimSpace(string(buf[0:n]))

		switch strings.ToLower(msg) {
		case "ls":
			{
				handleLS(&user, conn)
			}
		case "cd":
			{

			}
		case "pwd":
			{

			}
		case "exit":
			{
				return
			}
		default:
			{
				errMsg := ErrorMsg{Error: "Invalid command..."}
				data, _ := json.Marshal(errMsg)
				conn.Write(data)
			}
		}
	}
}

func handleError(err error, conn net.Conn) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	conn.Write([]byte(msg))

	return true
}

func handleLS(user *User, conn net.Conn) bool {
	wd, err := os.Open(user.WorkingDir)
	if handleError(err, conn) {
		return false
	}

	list, err := wd.Readdirnames(-1)
	if handleError(err, conn) {
		return false
	}

	dir := DirList{Directories: list}

	data, err := json.Marshal(dir)
	if handleError(err, conn) {
		return false
	}

	conn.Write(data)

	return true
}
