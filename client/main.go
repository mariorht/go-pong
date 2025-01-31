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
	width  = 80
	height = 24
	bgColor = "\033[48;5;235m" // Background color
	resetColor = "\033[0m"     // Reset color
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func drawGame(paddle1Y, paddle2Y, ballX, ballY int) {
	clearScreen()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Print(bgColor) // Set background color
			if y == 0 || y == height-1 {
				fmt.Print("-") // Top and bottom borders
			} else if x == 0 || x == width-1 {
				fmt.Print("|") // Left and right borders
			} else if x == 2 && y == paddle1Y {
				fmt.Print("|") // Paleta 1
			} else if x == width-3 && y == paddle2Y {
				fmt.Print("|") // Paleta 2
			} else if x == ballX && y == ballY {
				fmt.Print("O") // Pelota
			} else {
				fmt.Print(" ")
			}
			fmt.Print(resetColor) // Reset color
		}
		fmt.Println()
	}
}

func setRawMode() {
	var termios syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
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
