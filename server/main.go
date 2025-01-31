package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type GameState struct {
	Paddle1Y float64
	Paddle2Y float64
	BallX    float64
	BallY    float64
	BallDX   float64
	BallDY   float64
}

const paddleHeight = 3

var (
	state   = GameState{}
	clients = make(map[net.Conn]int)
	mu      sync.Mutex
)

func handleClient(conn net.Conn, player int) {
	defer conn.Close()
	fmt.Println("Player", player, "connected")

	// Send player number to client
	fmt.Fprintf(conn, "%d\n", player)

	if player == 2 {
		mu.Lock()
		state.BallX = 40
		state.BallY = 12
		state.BallDX = 1
		state.BallDY = 1
		mu.Unlock()
	}

	buf := make([]byte, 1)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Player", player, "disconnected")
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}
		fmt.Println("Player", player, "pressed", string(buf))
		mu.Lock()
		if player == 1 && buf[0] == 'w' {
			if state.Paddle1Y > 0 {
				state.Paddle1Y -= 1
			}
		} else if player == 1 && buf[0] == 's' {
			if state.Paddle1Y < 23-paddleHeight {
				state.Paddle1Y += 1
			}
		} else if player == 2 && buf[0] == 'w' {
			if state.Paddle2Y > 0 {
				state.Paddle2Y -= 1
			}
		} else if player == 2 && buf[0] == 's' {
			if state.Paddle2Y < 23-paddleHeight {
				state.Paddle2Y += 1
			}
		}
		mu.Unlock()
	}
}

func gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // ~10 FPS
	for range ticker.C {
		mu.Lock()

		// Mover la bola
		state.BallX += state.BallDX
		state.BallY += state.BallDY

		// Rebote en bordes verticales
		if state.BallY <= 0 {
			state.BallY = 0
			state.BallDY = -state.BallDY
		} else if state.BallY >= 23 {
			state.BallY = 23
			state.BallDY = -state.BallDY
		}

		// Rebote en bordes horizontales y paletas
		if state.BallX <= 3 {
			if state.BallY >= state.Paddle1Y && state.BallY <= state.Paddle1Y+paddleHeight {
				state.BallDX = -state.BallDX
			} else {
				state.BallX = 40
				state.BallY = 12
				state.BallDX = 1
				state.BallDY = 1
			}
		} else if state.BallX >= 76 {
			if state.BallY >= state.Paddle2Y && state.BallY <= state.Paddle2Y+paddleHeight {
				state.BallDX = -state.BallDX
			} else {
				state.BallX = 40
				state.BallY = 12
				state.BallDX = -1
				state.BallDY = 1
			}
		}

		// Enviar estado del juego a los clientes
		for conn := range clients {
			fmt.Fprintf(conn, "%d %d %d %d\n", int(state.Paddle1Y), int(state.Paddle2Y), int(state.BallX), int(state.BallY))
		}

		mu.Unlock()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}
	port := os.Args[1]

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error al iniciar servidor:", err)
		return
	}
	defer listener.Close()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error al obtener la IP:", err)
		return
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Servidor iniciado en la IP", ipnet.IP.String(), "y puerto", port)
			}
		}
	}

	go gameLoop()

	player := 1
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error en conexión:", err)
			continue
		}

		mu.Lock()
		clients[conn] = player
		mu.Unlock()

		go handleClient(conn, player)

		player++
		if player > 2 {
			player = 1
		}
	}
}
