package main

import (
	"encoding/json"
	"fmt"
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
	gameBoard = GameBoard{Board: [][]string{
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
		
		//TODO read json and unmarshal
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
				// TODO read json and unmarshal
				n, err = conn.Read(buf[:])
				if err != nil {
					fmt.Println(err)
					return
				}

				// TODO Parse and convert the data from Packet
				data := strings.TrimSpace(string(buf[0:n]))
				rowcol := strings.Split(data, ",")
				row, _ := strconv.Atoi(rowcol[0])
				col, _ := strconv.Atoi(rowcol[1])

				// Record player selection
				gameBoard.Board[row][col] = currentPlayer.Assigned

				// TODO check for winning state

				// Make current player state packet
				cpckt := Packet{Keys: make([]string, 1), Values: make([]string, 1)}

				cpckt.Keys[0] = PLAYER_STATE
				cpckt.Values[0] = "0"

				// Make waiting player state packet
				wpckt := Packet{Keys: make([]string, 3), Values: make([]string, 3)}

				wpckt.Keys[0] = BOARD_UPDATE
				wpckt.Values[0] = data

				wpckt.Keys[1] = BANNER
				wpckt.Values[1] = fmt.Sprintf("%s selected row: %d col: %d", currentPlayer.Name, row, col)

				wpckt.Keys[2] = PLAYER_STATE
				wpckt.Values[2] = "1"

				// Send json packets
				jcpckt, _ := json.Marshal(cpckt)
				jwpckt, _ := json.Marshal(wpckt)

				connections[currentPlayer.Name].Write(jcpckt)
				connections[waitingPlayer.Name].Write(jwpckt)

				// Update player states
				temp := currentPlayer
				currentPlayer = waitingPlayer
				waitingPlayer = temp
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

	// Create a packet of data to send
	cpckt := Packet{Keys: make([]string, 3), Values: make([]string, 3)}
	wpckt := Packet{Keys: make([]string, 3), Values: make([]string, 3)}

	// Send out banner message
	banner := fmt.Sprintf("%s is first move", currentPlayer.Name)

	cpckt.Keys[0] = BANNER
	cpckt.Values[0] = banner

	wpckt.Keys[0] = BANNER
	wpckt.Values[0] = banner

	// Update current player state
	cpckt.Keys[1] = ASSIGN
	cpckt.Values[1] = currentPlayer.Assigned

	cpckt.Keys[2] = PLAYER_STATE
	cpckt.Values[2] = "1"

	// Update waiting player state
	wpckt.Keys[1] = ASSIGN
	wpckt.Values[1] = waitingPlayer.Assigned

	wpckt.Keys[2] = PLAYER_STATE
	wpckt.Values[2] = "0"

	// Write the packets
	jcpckt, _ := json.Marshal(cpckt)
	jwpckt, _ := json.Marshal(wpckt)

	connections[currentPlayer.Name].Write(jcpckt)
	connections[waitingPlayer.Name].Write(jwpckt)
}

func cleanup(listener net.Listener) {
	for _, c := range connections {
		c.Close()
	}

	listener.Close()
}
