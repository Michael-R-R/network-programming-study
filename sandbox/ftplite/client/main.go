package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type DirList struct {
	Directories []string
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1200")
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		cmd = strings.TrimSpace(cmd)
		cmd = strings.ToLower(cmd)

		switch cmd {
		case "ls":
			{
				conn.Write([]byte(cmd))

				var buf [512]byte
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Println(err)
					continue
				}

				var dl DirList
				json.Unmarshal(buf[0:n], &dl)

				for _, d := range dl.Directories {
					fmt.Println(d)
				}
			}
		case "exit":
			{
				conn.Write([]byte(cmd))
				conn.Close()
				return
			}
		default:
			{

			}
		}
	}
}
