package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"strings"
)

type GameBoard struct {
	Board [][]string
}

type Player struct {
	Id       int
	Name     string
	Assigned string
	Conn     net.Conn
}

type Packet struct {
	Keys   []string
	Values []string
}

const ( // Board piece types
	NONE    = ""
	X_PIECE = "X"
	O_PIECE = "O"
)

const ( // Message types
	BOARD_UPDATE = "10"
	PLAYER_STATE = "11"
	BANNER       = "12"
)

var gameBoard GameBoard
var player1 Player
var player2 Player

var currentPlayer *Player
var waitingPlayer *Player

func init() {
	gameBoard = GameBoard{Board: [][]string{
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
	}}
	player1 = Player{Id: 1, Name: "Player1", Assigned: NONE, Conn: nil}
	player2 = Player{Id: 2, Name: "Player2", Assigned: NONE, Conn: nil}
	currentPlayer = nil
	waitingPlayer = nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <host> <port>\n", os.Args[0])
	}

	host := os.Args[1]
	port := os.Args[2]
	_, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Fatalf("Invalid host: %s", host)
	}

	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalln(err)
	}

	initConnections(listener)
	gameLoop()
	cleanup(listener)
}

func initConnections(listener net.Listener) {
	count := 0

	for count < 2 {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Connected: %s\n", conn.RemoteAddr().String())

		if count == 0 {
			player1.Conn = conn
		} else {
			player2.Conn = conn
		}

		count++
	}
}

func gameLoop() {
	initPieces()

	for {
		// Read current client data
		data, err := readAll(currentPlayer.Conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var cpckt Packet
		_ = json.Unmarshal(data, &cpckt)

		for i, k := range cpckt.Keys {
			switch k {
			case BOARD_UPDATE:
				{
					// Parse and convert the row/col data
					value := cpckt.Values[i]
					rowcol := strings.Split(value, ",")
					row, _ := strconv.Atoi(rowcol[0])
					col, _ := strconv.Atoi(rowcol[1])

					// TODO verify selection

					// Record player selection
					gameBoard.Board[row][col] = currentPlayer.Assigned

					// TODO check for winning state

					// Make current player packet
					cpckt := Packet{Keys: make([]string, 2), Values: make([]string, 2)}

					cpckt.Keys[0] = BOARD_UPDATE
					cpckt.Values[0] = currentPlayer.Assigned + "," + value

					cpckt.Keys[1] = PLAYER_STATE
					cpckt.Values[1] = "0"

					// Make waiting player packet
					wpckt := Packet{Keys: make([]string, 3), Values: make([]string, 3)}

					wpckt.Keys[0] = BOARD_UPDATE
					wpckt.Values[0] = currentPlayer.Assigned + "," + value

					wpckt.Keys[1] = BANNER
					wpckt.Values[1] = fmt.Sprintf("%s selected row: %d col: %d", currentPlayer.Name, row, col)

					wpckt.Keys[2] = PLAYER_STATE
					wpckt.Values[2] = "1"

					// Send json packets
					jcpckt, _ := json.Marshal(cpckt)
					jwpckt, _ := json.Marshal(wpckt)

					currentPlayer.Conn.Write(jcpckt)
					waitingPlayer.Conn.Write(jwpckt)

					// Update player states
					temp := currentPlayer
					currentPlayer = waitingPlayer
					waitingPlayer = temp
				}
			}
		}
	}
}

func initPieces() {
	// Randomly assign pieces
	// X = first move
	n := 1 + rand.Int()*(100-1)
	if n < 50 {
		player1.Assigned = X_PIECE
		player2.Assigned = O_PIECE

		currentPlayer = &player1
		waitingPlayer = &player2
	} else {
		player1.Assigned = O_PIECE
		player2.Assigned = X_PIECE

		currentPlayer = &player2
		waitingPlayer = &player1
	}

	banner := fmt.Sprintf("%s is first move", currentPlayer.Name)

	// Create current player packet
	cpckt := Packet{Keys: make([]string, 2), Values: make([]string, 2)}

	cpckt.Keys[0] = BANNER
	cpckt.Values[0] = banner

	cpckt.Keys[1] = PLAYER_STATE
	cpckt.Values[1] = "1"

	// Create waiting player packet
	wpckt := Packet{Keys: make([]string, 2), Values: make([]string, 2)}

	wpckt.Keys[0] = BANNER
	wpckt.Values[0] = banner

	wpckt.Keys[1] = PLAYER_STATE
	wpckt.Values[1] = "0"

	// Write the packets
	jcpckt, _ := json.Marshal(cpckt)
	jwpckt, _ := json.Marshal(wpckt)

	currentPlayer.Conn.Write(jcpckt)
	waitingPlayer.Conn.Write(jwpckt)
}

func readAll(conn net.Conn) ([]byte, error) {
	var buf [512]byte
	result := bytes.NewBuffer(nil)

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

func cleanup(listener net.Listener) {
	player1.Conn.Close()
	player2.Conn.Close()

	listener.Close()
}
