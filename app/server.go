package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory = flag.String("directory", "/", "define root directory")

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(conn net.Conn) {
			defer func() {
				err = conn.Close()
				if err != nil {
					fmt.Println("Error closing connection: ", err.Error())
					os.Exit(1)
				}
			}()

			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading data from the connection: ", err.Error())
				os.Exit(1)
			}

			parts := strings.Split(string(buffer), "\r\n")

			request := strings.Split(parts[0], " ")

			route := request[1]

			path := strings.Split(route, "/")[1:]

			var response string
			switch path[0] {
			case "":
				response = "HTTP/1.1 200 OK\r\n\r\n"
			case "echo":
				if len(path) <= 1 {
					response = "HTTP/1.1 400 Bad Request\r\n\r\n"
					break
				}
				response = fmt.Sprintf(
					"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
					len(path[1]),
					path[1],
				)
			case "files":
				if directory == nil {
					response = "HTTP/1.1 404 Not Found\r\n\r\n"
					break
				}

				if len(path) <= 1 {
					response = "HTTP/1.1 400 Bad Request\r\n\r\n"
					break
				}

				filename := path[1]

				if (*directory)[len(*directory)-1] != '/' {
					*directory += "/"
				}

				content, err := os.ReadFile(*directory + filename)
				if err != nil {
					response = "HTTP/1.1 404 Not Found\r\n\r\n"
					break
				}

				response = fmt.Sprintf(
					"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
					len(content), content,
				)
			case "user-agent":
				// iterate over the headers until finding an empty part (between the last \r\n (empty) \r\n)
				curHeaderIdx := 1
				for curHeaderIdx < len(parts)-1 && parts[curHeaderIdx] != "" {
					kv := strings.Split(parts[curHeaderIdx], ": ")

					if kv[0] == "User-Agent" {
						response = fmt.Sprintf(
							"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
							len(kv[1]), kv[1],
						)
						break
					}

					curHeaderIdx += 1
				}
			default:
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			}

			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to the connection: ", err.Error())
				os.Exit(1)
			}
		}(conn)
	}
}
