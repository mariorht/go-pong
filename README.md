# Go Pong

Go Pong is a simple terminal-based Pong game implemented in Go. The game consists of a server and multiple clients. Players can connect to the server and control their paddles to play the game.

## How to Play

### Server

1. Start the server by running the following command:
    ```sh
    go run server/main.go <port>
    ```
    Replace `<port>` with the port number you want the server to listen on.

### Client

1. Start the client by running the following command:
    ```sh
    go run client/main.go <server_ip> <server_port>
    ```
    Replace `<server_ip>` with the IP address of the server and `<server_port>` with the port number the server is listening on.

### Controls

- Player 1:
  - `w`: Move paddle up
  - `s`: Move paddle down

- Player 2:
  - `w`: Move paddle up
  - `s`: Move paddle down

### Game Rules

- The game starts when two players are connected to the server.
- Each player controls a paddle on their side of the screen.
- The objective is to hit the ball with your paddle and prevent it from going past your paddle.
- The ball will bounce off the paddles and the top and bottom borders of the screen.
- If the ball goes past a player's paddle, the ball will be reset to the center of the screen.

## Project Structure

- `server/main.go`: The server code that handles client connections and game logic.
- `client/main.go`: The client code that connects to the server and renders the game.

## Dependencies

- Go 1.16 or later

## License

This project is licensed under the MIT License.