package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var _ = net.Listen
var _ = os.Exit
var file_path = "files"

func main() {
	fmt.Println("Logs from your program will appear here!")
	fmt.Println("-----------------------------------------")

	args := os.Args
	if len(args) > 1 {
		if args[1] == "--directory" {
			file_path = args[2]
		}
		fmt.Printf("%s using this for file directory\n", file_path)
	}

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

	// this is pure byte for the request
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

		var checker bool = false
		multienco := strings.Split(headers["accept-encoding"], ", ")
		fmt.Printf("this is the multienco %s", multienco)
		if len(multienco) > 1 {
			for _, v := range multienco {
				if v == "gzip" {
					checker = true
				}
			}
		}

		if headers["accept-encoding"] == "gzip" {
			checker = true
		}

		if checker {
			val, _ := GzipCompress(echoed)
			response := fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s",
				len(val), val)
			conn.Write([]byte(response))
			return
		}

		writeResponse(conn, 200, echoed)
		return
	}

	if method == "GET" && strings.HasPrefix(path, "/files/") {
		file_name := strings.TrimPrefix(path, "/files/")
		if checkfileexist(file_path + file_name) {
			file_content, _ := os.ReadFile(file_path + file_name)
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

	if method == "POST" && strings.HasPrefix(path, "/files/") {
		buffer := make([]byte, reader.Size())
		_, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}

		content_length := headers["content-length"]
		length, err := strconv.Atoi(content_length)
		file_name := strings.TrimPrefix(path, "/files/")
		err = os.WriteFile(file_path+file_name, buffer[:length], 0644)
		if err != nil {
			panic(err)
		}
		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		return
	}

	writeResponse(conn, 404, " ")
}
func GzipCompress(str string) (string, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(str))
	if err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func checkfileexist(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
