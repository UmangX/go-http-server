package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

//HTTP/1.1 404 Not Found\r\n\r\n

func main() {
	fmt.Println("Logs from your program will appear here!")
	line, _ := net.Listen("tcp", ":4221")
	for {
		conn, _ := line.Accept()
		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()

	// this is the buffer for the data that the client is sending
	buffer := make([]byte, 2048)
	conn.Read(buffer)

	content := string(buffer)
	lines := strings.Split(content, "\n")
	// lines has all the line or something else
	fmt.Println("request from the client")
	for _, val := range lines {
		fmt.Println(val)
	}

	//this is some fine work here
	// when this strings is done will work on using the direct bytes for this
	req_type := strings.Split(lines[0], " ")[0]
	url_text := strings.Split(lines[0], " ")[1]
	fmt.Printf("req type %v  url_text %v", req_type, url_text)
	// this 'response' is a plain and simple string that is send when req is valid

	response200 := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, world!"

	response404 := "HTTP/1.1 404 Not Found\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, world!"

	if url_text == "/" {
		conn.Write([]byte(response200))
	} else {
		conn.Write([]byte(response404))
	}

}
