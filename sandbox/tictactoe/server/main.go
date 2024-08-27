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
	PLAYER_STATE = "11"
	BANNER       = "12"
	ASSIGN       = "13"
)

var connections map[string]net.Conn // [player id, connection]
var gameBoard GameBoard
var player1 Player
var player2 Player

var currentPlayer *Player
var waitingPlayer *Player

func init() {
	connections = make(map[string]net.Conn)
	gameBoard = GameBoard{Board: [][]int{
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
		{NONE, NONE, NONE},
	}}
	player1 = Player{Id: 1, Name: "Player1", Assigned: NONE}
	player2 = Player{Id: 2, Name: "Player2", Assigned: NONE}
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
	for len(connections) < 2 {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Connected: %s\n", conn.RemoteAddr().String())

		if len(connections) == 0 {
			connections[player1.Name] = conn
		} else {
			connections[player2.Name] = conn
		}
	}
}

func gameLoop() {
	initPieces()

	for {
		conn := connections[currentPlayer.Name]

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

				// Parse and convert the data
				data := strings.TrimSpace(string(buf[0:n]))
				rowcol := strings.Split(data, ",")
				row, _ := strconv.Atoi(rowcol[0])
				col, _ := strconv.Atoi(rowcol[1])

				// Record player selection
				gameBoard.Board[row][col] = currentPlayer.Assigned

				// TODO check for winning state

				// Update waiting player board state
				connections[waitingPlayer.Name].Write([]byte(BOARD_UPDATE))
				connections[waitingPlayer.Name].Write([]byte(data))

				// Update waiting player banner
				banner := fmt.Sprintf("%s selected row: %d col: %d", currentPlayer.Name, row, col)
				connections[waitingPlayer.Name].Write([]byte(BANNER))
				connections[waitingPlayer.Name].Write([]byte(banner))

				// Update player states
				temp := currentPlayer
				currentPlayer = waitingPlayer
				waitingPlayer = temp

				connections[currentPlayer.Name].Write([]byte(PLAYER_STATE))
				connections[currentPlayer.Name].Write([]byte("1"))

				connections[waitingPlayer.Name].Write([]byte(PLAYER_STATE))
				connections[waitingPlayer.Name].Write([]byte("0"))
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

	// Update connections about assigned pieces
	banner := fmt.Sprintf("%s is first move", currentPlayer.Name)
	for _, c := range connections {
		c.Write([]byte(BANNER))
		c.Write([]byte(banner))
	}

	// Update current player turn state
	conn1 := connections[currentPlayer.Name]
	conn1.Write([]byte(ASSIGN))
	conn1.Write([]byte(strconv.Itoa(currentPlayer.Assigned)))
	conn1.Write([]byte(PLAYER_STATE))
	conn1.Write([]byte("1"))

	// Update waiting player turn state
	conn2 := connections[waitingPlayer.Name]
	conn2.Write([]byte(ASSIGN))
	conn2.Write([]byte(strconv.Itoa(currentPlayer.Assigned)))
	conn2.Write([]byte(PLAYER_STATE))
	conn2.Write([]byte("0"))
}

func cleanup(listener net.Listener) {
	for _, c := range connections {
		c.Close()
	}

	listener.Close()
}
