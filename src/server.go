package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"redis/src/lexer"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error listening on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}

		answer := lexer.ReadManager(buf[:n])

		_, err = conn.Write([]byte(answer))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}
