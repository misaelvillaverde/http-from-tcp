package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
			os.Exit(1)
		}
	}()

	request := make([]byte, 50)
	_, err = conn.Read(request)
	if err != nil {
		fmt.Println("Error reading data from the connection: ", err.Error())
		os.Exit(1)
	}

	parts := strings.Split(string(request), "\r\n")

	requestLine := strings.Split(parts[0], " ")

	switch requestLine[1] {
	case "/":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Error writing data to the connection: ", err.Error())
		os.Exit(1)
	}
}
