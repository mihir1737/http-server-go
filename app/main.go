package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Accepts the connection
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

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
	}

	// construct the row HTTP response string
	// Each line end with \r\n (CRLF)

	var response []byte

	if request.RequestURI == "/index.html" || request.RequestURI == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	// Write the response to the connection
	_, err = conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}

}
