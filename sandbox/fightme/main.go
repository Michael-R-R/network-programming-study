package main

import (
	"fmt"
	"math/rand/v2"
	"net"
)

const (
	WAITING  = iota
	GAMEOVER = iota
	PLAYER1  = iota
	PLAYER2  = iota
)

type Player struct {
	health  float64
	defense float64
}

var state uint
var player1 Player
var player2 Player
var connections map[int]net.Conn

func init() {
	state = WAITING

	player1 = Player{
		health:  100.0,
		defense: 1.0 + rand.Float64()*(10.0-1.0),
	}

	player2 = Player{
		health:  100.0,
		defense: 1.0 + rand.Float64()*(10.0-1.0),
	}

	connections = make(map[int]net.Conn)
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1200")
	if err != nil {
		fmt.Println(err)
		return
	}

	count := 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if count < 2 {
			count++
			connections[count] = conn

			if count == 2 {
				state = PLAYER1

				handleConnections()
				listener.Close()
				return
			}
		}
	}
}

func handleConnections() {
	for {
		switch state {
		case PLAYER1:
			{
				result := handlePlayer(connections[1], &player2)
				announce("Player1", "Player2", result, player2.health)
				state = PLAYER2
				if player2.health <= 0.0 {
					state = GAMEOVER
				}
			}
		case PLAYER2:
			{
				result := handlePlayer(connections[2], &player1)
				announce("Player2", "Player1", result, player1.health)
				state = PLAYER1
				if player1.health <= 0.0 {
					state = GAMEOVER
				}
			}
		case WAITING:
			{
				// do nothing
			}
		case GAMEOVER:
			{
				for _, conn := range connections {
					conn.Write([]byte("Game Over!"))
					conn.Close()
				}
				return
			}
		}
	}
}

func handlePlayer(conn net.Conn, def *Player) (dmg float64) {
	conn.Write([]byte("PRESS ENTER\n"))

	var n int
	var buf [512]byte
	for string(buf[:n]) != "\r\n" {
		n, _ = conn.Read(buf[0:])
	}

	dmg = (25.0 + rand.Float64()*(50.0-25.0)) - def.defense
	if dmg < 0.0 {
		dmg = 0.0
	}

	def.health -= dmg

	return dmg
}

func announce(attker, defender string, dmgdone, remaining float64) {
	msg := fmt.Sprintf(
		"%s did %.2f to %s - %s HP remaining = %.2f\n",
		attker, dmgdone, defender, defender, remaining,
	)

	for _, conn := range connections {
		conn.Write([]byte(msg))
	}
}
