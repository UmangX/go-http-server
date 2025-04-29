package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var filePath = "files"

func main() {
	fmt.Println("Logs from your program will appear here!")
	fmt.Println("-----------------------------------------")

	args := os.Args
	if len(args) > 2 && args[1] == "--directory" {
		filePath = args[2]
	}
	fmt.Printf("Using directory: %s\n", filePath)

	listener, err := net.Listen("tcp", ":4221")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(noTimeout())

		requestLine, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading request:", err)
			}
			return
		}
		requestLine = strings.TrimSpace(requestLine)
		parts := strings.Split(requestLine, " ")
		if len(parts) < 2 {
			writeResponse(conn, 400, "Bad Request", nil)
			return
		}
		method, path := parts[0], parts[1]

		// Parse headers
		headers := make(map[string]string)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Header read error:", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headers[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
			}
		}

		// Handle routes
		switch {
		case method == "GET" && path == "/":
			writeResponse(conn, 200, "OK", []byte("hello"))

		case method == "GET" && strings.HasPrefix(path, "/echo/"):
			message := strings.TrimPrefix(path, "/echo/")
			if acceptsGzip(headers["accept-encoding"]) {
				body, _ := gzipCompress([]byte(message))
				writeRawResponse(conn, 200, "text/plain", body, map[string]string{"Content-Encoding": "gzip"})
			} else {
				writeResponse(conn, 200, "OK", []byte(message))
			}

		case method == "GET" && strings.HasPrefix(path, "/files/"):
			filename := strings.TrimPrefix(path, "/files/")
			fullPath := filepath.Join(filePath, filename)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				writeResponse(conn, 404, "Not Found", nil)
			} else {
				writeRawResponse(conn, 200, "application/octet-stream", data, nil)
			}

		case method == "POST" && strings.HasPrefix(path, "/files/"):
			filename := strings.TrimPrefix(path, "/files/")
			fullPath := filepath.Join(filePath, filename)
			length, _ := strconv.Atoi(headers["content-length"])
			body := make([]byte, length)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				fmt.Println("Error reading POST body:", err)
				writeResponse(conn, 500, "Internal Server Error", nil)
				return
			}
			err = os.WriteFile(fullPath, body, 0644)
			if err != nil {
				writeResponse(conn, 500, "Internal Server Error", nil)
			} else {
				writeResponse(conn, 201, "Created", nil)
			}

		case method == "GET" && path == "/user-agent":
			writeResponse(conn, 200, "OK", []byte(headers["user-agent"]))

		default:
			writeResponse(conn, 404, "Not Found", nil)
		}
	}
}

// gzipCompress compresses a byte slice into gzip format
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

// writeResponse writes a plain text HTTP response
func writeResponse(conn net.Conn, status int, statusText string, body []byte) {
	if body == nil {
		body = []byte{}
	}
	writeRawResponse(conn, status, "text/plain", body, nil)
}

// writeRawResponse writes a full HTTP response with optional headers
func writeRawResponse(conn net.Conn, status int, contentType string, body []byte, extraHeaders map[string]string) {
	headers := fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, statusText(status))
	headers += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	headers += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	headers += "Connection: keep-alive\r\n"

	for k, v := range extraHeaders {
		headers += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headers += "\r\n"

	conn.Write([]byte(headers))
	conn.Write(body)
}

func acceptsGzip(header string) bool {
	for _, val := range strings.Split(header, ",") {
		if strings.TrimSpace(val) == "gzip" {
			return true
		}
	}
	return false
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

func noTimeout() (t time.Time) {
	// disables timeout by returning zero time
	return
}
