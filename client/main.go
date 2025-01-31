package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	width       = 80
	height      = 24
	bgColor     = "\033[48;5;235m" // Background color
	resetColor  = "\033[0m"        // Reset color
	paddleHeight = 3
)

var player int

func moveToTopLeft() {
	fmt.Print("\033[H")
}

func drawGame(paddle1Y, paddle2Y, ballX, ballY int) {
	var buffer strings.Builder

	buffer.WriteString(fmt.Sprintf("Player %d\n", player)) // Display player number

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			buffer.WriteString(bgColor) // Set background color
			if y == 0 || y == height-1 {
				buffer.WriteString("-") // Top and bottom borders
			} else if x == 0 || x == width-1 {
				buffer.WriteString("|") // Left and right borders
			} else if x == 2 && y >= paddle1Y && y < paddle1Y+paddleHeight {
				buffer.WriteString("|") // Paleta 1
			} else if x == width-3 && y >= paddle2Y && y < paddle2Y+paddleHeight {
				buffer.WriteString("|") // Paleta 2
			} else if x == ballX && y == ballY {
				buffer.WriteString("O") // Pelota
			} else {
				buffer.WriteString(" ")
			}
			buffer.WriteString(resetColor) // Reset color
		}
		buffer.WriteString("\n")
	}

	moveToTopLeft()
	fmt.Print(buffer.String())
}

func setRawMode() {
	var termios syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <server_ip> <server_port>")
		return
	}
	serverIP := os.Args[1]
	serverPort := os.Args[2]

	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		log.Fatal("No se pudo conectar al servidor:", err)
	}
	defer conn.Close()

	setRawMode()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadByte()
			conn.Write([]byte{input}) // Enviar solo la primera letra ("w" o "s")
		}
	}()

	// Receive player number from server
	playerMessage, _ := bufio.NewReader(conn).ReadString('\n')
	player, _ = strconv.Atoi(strings.TrimSpace(playerMessage))

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		data := strings.Fields(message)
		if len(data) != 4 {
			continue
		}

		paddle1Y, _ := strconv.Atoi(data[0])
		paddle2Y, _ := strconv.Atoi(data[1])
		ballX, _ := strconv.Atoi(data[2])
		ballY, _ := strconv.Atoi(data[3])

		drawGame(paddle1Y, paddle2Y, ballX, ballY)

		time.Sleep(100 * time.Millisecond) // Refresco de 10 FPS
	}
}
