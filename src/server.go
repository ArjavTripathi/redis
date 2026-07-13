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

	PORT := os.Getenv("PORT")
	fmt.Printf("Listeneing on port %s", PORT)

	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Printf("Error listening on port %s: %v", PORT, err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v", err)
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading from connection: %v", err)
			os.Exit(1)
		}
		answer := lexer.ReadStream(buf)
		conn.Write([]byte(answer))
	}
}
