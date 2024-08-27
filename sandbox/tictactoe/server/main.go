package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"strings"
)

type GameBoard struct {
	Board [][]int
}

type Player struct {
	Id       int
	Name     string
	Assigned int
}

const ( // Board piece types
	NONE    = 0
	X_PIECE = 1
	O_PIECE = 2
)

const ( // Message types
	BOARD_UPDATE = "10"
	PLAYER_TURN  = "11"
	BANNER       = "12"
	ASSIGN       = "13"
)

var connections map[Player]net.Conn // [player id, connection]
var gameBoard GameBoard
var player1 Player
var player2 Player
var currentPlayer *Player

func init() {
	connections = make(map[Player]net.Conn)
	gameBoard = GameBoard{Board: [][]int{
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
	}}
	player1 = Player{Id: 1, Name: "Player1", Assigned: NONE}
	player2 = Player{Id: 2, Name: "Player2", Assigned: NONE}
	currentPlayer = nil
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
	for len(connections) < 2 {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(connections) == 0 {
			connections[player1] = conn
		} else {
			connections[player2] = conn
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
	} else {
		player1.Assigned = O_PIECE
		player2.Assigned = X_PIECE
		currentPlayer = &player2
	}

	// Update connections about assigned pieces
	msg := fmt.Sprintf("%s is first move", currentPlayer.Name)
	for p, c := range connections {
		c.Write([]byte(ASSIGN))
		c.Write([]byte(strconv.Itoa(p.Assigned)))

		c.Write([]byte(BANNER))
		c.Write([]byte(msg))
	}

	// Set first move current player
	connections[*currentPlayer].Write([]byte(PLAYER_TURN))
	connections[*currentPlayer].Write([]byte("1"))
}

func gameLoop() {
	initPieces()

	for {
		conn := connections[*currentPlayer]

		var buf [512]byte

		// Read message type
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := string(buf[0:n])
		msg = strings.TrimSpace(msg)

		switch msg {
		case BOARD_UPDATE:
			{
				n, err = conn.Read(buf[:])
				if err != nil {
					fmt.Println(err)
					return
				}

				data := strings.TrimSpace(string(buf[0:n]))
				rowcol := strings.Split(data, ",")
				row, _ := strconv.Atoi(rowcol[0])
				col, _ := strconv.Atoi(rowcol[1])

				// Record selection

				// Update players with banner

				// Set player states

				// Assign new current player
			}
		}
	}
}

func cleanup(listener net.Listener) {
	for _, c := range connections {
		c.Close()
	}

	listener.Close()
}
