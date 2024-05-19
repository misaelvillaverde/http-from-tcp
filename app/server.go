package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
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

			body := parts[len(parts)-1]

			request := strings.Split(parts[0], " ")

			method := request[0] // GET, POST

			route := request[1]

			path := strings.Split(route, "/")[1:]

			var response string
			switch path[0] {
			case "":
				response = "HTTP/1.1 200 OK" + CLRF + CLRF
			case "echo":
				if len(path) <= 1 {
					response = "HTTP/1.1 400 Bad Request" + CLRF + CLRF
					break
				}

				acceptEncoding, _ := getHeader(parts, "Accept-Encoding")

				headers := NewHeaders().
					setEncoding(acceptEncoding).
					setContentType("text/plain").
					setContentLength(len(path[1]))

				response =
					"HTTP/1.1 200 OK" +
						CLRF +
						headers.String() +
						CLRF +
						path[1]
			case "files":
				if directory == nil {
					response = "HTTP/1.1 404 Not Found" + CLRF + CLRF
					break
				}

				if len(path) <= 1 {
					response = "HTTP/1.1 400 Bad Request" + CLRF + CLRF
					break
				}

				filename := path[1]

				if (*directory)[len(*directory)-1] != '/' {
					*directory += "/"
				}

				switch method {
				case "GET":
					content, err := os.ReadFile(*directory + filename)
					if err != nil {
						response = "HTTP/1.1 404 Not Found" + CLRF + CLRF
						break
					}

					headers := NewHeaders().
						setContentType("application/octet-stream").
						setContentLength(len(content))

					response =
						"HTTP/1.1 200 OK" +
							CLRF +
							headers.String() +
							CLRF +
							string(content)
				case "POST":
					contentLength, _ := getHeader(parts, "Content-Length")
					length, _ := strconv.Atoi(contentLength)

					body = body[:length]

					err = os.WriteFile(*directory+filename, []byte(body), 0644)
					if err != nil {
						response = "HTTP/1.1 400 Bad Request" + CLRF + CLRF + err.Error()
						break
					}

					response = "HTTP/1.1 201 Created" + CLRF + CLRF
				}
			case "user-agent":
				userAgent, err := getHeader(parts, "User-Agent")
				if err != nil {
					response = "HTTP/1.1 404 Not Found" + CLRF + CLRF + err.Error()
					break
				}

				headers := NewHeaders().
					setContentType("text/plain").
					setContentLength(len(userAgent))

				response =
					"HTTP/1.1 200 OK" +
						CLRF +
						headers.String() +
						CLRF +
						userAgent
			default:
				response = "HTTP/1.1 404 Not Found" + CLRF + CLRF
			}

			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to the connection: ", err.Error())
				os.Exit(1)
			}
		}(conn)
	}
}

func getHeader(parts []string, key string) (string, error) {
	// iterate over the headers until finding an empty part (between the last \r\n (empty) \r\n)
	curHeaderIdx := 1
	for curHeaderIdx < len(parts)-1 && parts[curHeaderIdx] != "" {
		kv := strings.Split(parts[curHeaderIdx], ": ")

		if kv[0] == key {
			return kv[1], nil
		}

		curHeaderIdx += 1
	}
	return "", errors.New("Could not find header")
}
