package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")
	fmt.Println("-----------------------------------------")

	listener, _ := net.Listen("tcp", ":4221")
	for {
		conn, _ := listener.Accept()
		go handleConn(conn)
	}
}

func writeResponse(conn net.Conn, statusCode int, body string) {
	statusText := map[int]string{
		200: "OK",
		404: "Not Found",
	}[statusCode]

	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, statusText, len(body), body,
	)

	conn.Write([]byte(response))
}

func writeResponseforfile(conn net.Conn, body string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(response))
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// this is for the lines which are the requests
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		writeResponse(conn, 404, " ")
		return
	}

	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")
	fmt.Printf("handling the endpoint : %v\n", parts)
	if len(parts) < 2 {
		writeResponse(conn, 404, " ")
		return
	}

	method := parts[0]
	path := parts[1]

	// Read headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			key := strings.ToLower(strings.TrimSpace(headerParts[0]))
			value := strings.TrimSpace(headerParts[1])
			headers[key] = value
		}
	}

	// Handle requests
	if method == "GET" && path == "/" {
		writeResponse(conn, 200, "hello")
		return
	}

	if method == "GET" && strings.HasPrefix(path, "/echo/") {
		echoed := strings.TrimPrefix(path, "/echo/")
		writeResponse(conn, 200, echoed)
		return
	}

	if method == "GET" && strings.HasPrefix(path, "/files/") {
		file_name := strings.TrimPrefix(path, "/files/")
		if checkfileexist("./tmp/" + file_name) {
			file_content, _ := os.ReadFile("./tmp/" + file_name)
			writeResponseforfile(conn, string(file_content))
			return
		}
		writeResponse(conn, 404, " ")
		return
	}

	if method == "GET" && path == "/user-agent" {
		userAgent := headers["user-agent"]
		writeResponse(conn, 200, userAgent)
		return
	}

	writeResponse(conn, 404, " ")
}

func checkfileexist(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
