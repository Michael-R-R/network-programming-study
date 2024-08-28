package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

type User struct {
	WorkingDir string
	Conn       net.Conn
}

type Packet struct {
	Keys   []string
	Values []string
}

func init() {

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

		user := User{WorkingDir: ".", Conn: conn}

		go handleUser(user)
	}
}

func handleUser(user User) {
	defer user.Conn.Close()

	conn := user.Conn

	for {
		// Read client packet
		data, err := readAll(conn)
		if handleError(err, conn) {
			continue
		}

		var packet Packet
		err = json.Unmarshal(data, &packet)
		if handleError(err, conn) {
			continue
		}

		for i, key := range packet.Keys {
			key = strings.TrimSpace(key)
			key = strings.ToLower(key)
			value := packet.Values[i]

			switch key {
			case LS:
				{
					handleLS(&user)
				}
			case EXIT:
				{
					handleExit(&user)
					return
				}
			default:
				{
					err = fmt.Errorf("invalid command: %s %s", key, value)
					handleError(err, conn)
				}
			}
		}
	}
}

func handleError(err error, conn net.Conn) bool {
	if err == nil {
		return false
	}

	packet := Packet{Keys: make([]string, 1), Values: make([]string, 1)}
	packet.Keys[0] = ERROR
	packet.Values[0] = err.Error()

	data, _ := json.Marshal(packet)

	conn.Write(data)

	return true
}

func handleLS(user *User) bool {
	conn := user.Conn

	wd, err := os.Open(user.WorkingDir)
	if handleError(err, conn) {
		return false
	}

	dlist, err := wd.Readdirnames(-1)
	if handleError(err, conn) {
		return false
	}

	var val string
	for _, d := range dlist {
		val += fmt.Sprintf("%s,", d)
	}

	packet := Packet{Keys: make([]string, 1), Values: make([]string, 1)}
	packet.Keys[0] = LS
	packet.Values[0] = val

	data, _ := json.Marshal(packet)

	conn.Write(data)

	return true
}

func handleExit(user *User) {
	conn := user.Conn

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
