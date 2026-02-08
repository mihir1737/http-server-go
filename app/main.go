package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const supportedEncoding = "gzip"

func handleEcho(request http.Request, response []byte, path string) {
	echoStr := path[6:]

	clientAcceptEncoding := request.Header.Get("Accept-Encoding")

	if clientAcceptEncoding != "" {
		if clientAcceptEncoding == supportedEncoding {
			response = fmt.Appendf(
				nil,
				"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\n\r\n",
				clientAcceptEncoding,
			)
		} else {
			response = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")
		}
	} else {
		response = fmt.Appendf(
			nil,
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(echoStr),
			echoStr,
		)
	}
}

func handlePostFile(request http.Request, response []byte, directory string, path string) {
	contentLength := request.ContentLength
	contentBytes := make([]byte, contentLength)

	_, err := request.Body.Read(contentBytes)

	if err != nil {
		response = []byte("HTTP/1.1 500 Interal server Error\r\n\r\n")
	}

	filePath := directory + path[7:]
	err = os.WriteFile(filePath, contentBytes, 0644)

	if err != nil {
		response = []byte("HTTP/1.1 500 Interal server\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 201 Created\r\n\r\n\r\n\r\n%s")
	}
}

func handleGetFile(response []byte, directory string, path string) {
	filePath := directory + path[7:]
	content, err := os.ReadFile(filePath)

	if err != nil {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	} else {
		response = fmt.Appendf(
			nil,
			"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
			len(content),
			string(content),
		)
	}
}

func handleRequest(conn net.Conn, directory string) {
	// Always close your connections!
	// Used defer means, it works after the function returns
	// benifit of placing this statement here is in case of error as well
	// the connection will be closed.
	defer conn.Close()

	// wrap the connection in a buffered reader
	reader := bufio.NewReader(conn)

	// Parse the request
	request, err := http.ReadRequest(reader)

	if err != nil {
		fmt.Println("Error parsing the request.", err)
		return
	}
	// construct the row HTTP response string
	// Each line end with \r\n (CRLF)

	var response []byte

	path := request.URL.Path
	switch {
	case path == "/index.html" || path == "/":
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")

	case path == "/user-agent":
		response = fmt.Appendf(
			nil,
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(request.UserAgent()),
			request.UserAgent(),
		)

	case strings.HasPrefix(path, "/echo/"):
		handleEcho(*request, response, path)

	case strings.HasPrefix(path, "/files/") && request.Method == "POST":
		handlePostFile(*request, response, directory, path)

	case strings.HasPrefix(path, "/files/"):
		handleGetFile(response, directory, path)

	default:
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	// Write the response to the connection
	_, err = conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}
func main() {
	// To take the directory as input for getting files
	args := os.Args
	directory := ""

	if len(args) > 1 {
		directory = args[2]
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		// Accepts the connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn, directory)
	}
}
