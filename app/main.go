package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func generate_content_respone(content string) string {

	if content == " " {
		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: 0\r\n" +
			"\r\n" +
			""
		return response
	}

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: " + strconv.Itoa(len(content)) + "\r\n" +
		"\r\n" + content

	fmt.Printf(response)

	return response

}

func HandleConn(conn net.Conn) {
	defer conn.Close()

	// this is the buffer for the data that the client is sending
	buffer := make([]byte, 2048)
	conn.Read(buffer)

	content := string(buffer)

	// this line is the main content seperate through new line
	lines := strings.Split(content, "\n")

	// when this strings is done will work on using the direct bytes for this
	//req_type := strings.Split(lines[0], " ")[0]
	url_text := strings.Split(lines[0], " ")[1]

	url_lines := strings.Split(url_text, "/")
	if lines[1] == "echo" {
		//fmt.Printf("this is the echo link with /%v \n", lines[2])
		follow_up := url_lines[2]
		urlstr := url_lines[2][0:]
		if follow_up == "/" {
			conn.Write([]byte(generate_content_respone(" ")))
			return
		} else {
			conn.Write([]byte(generate_content_respone(urlstr)))
			return
		}
	}

	// write for the header actions
	if url_lines[1] == "user-agent" {
		user_agent := strings.Split(lines[3], ":")[1]
		fmt.Println(user_agent)
		conn.Write([]byte(generate_content_respone(user_agent)))
	}

	//fmt.Printf("req type %v  url_text %v len for the lines in url_text : %v  this is the first in lines : %v \n", req_type, url_text, len(lines), lines[0])
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

	switch link := url_text; link {
	case "/":
		conn.Write([]byte(response200))
	case "/echo": //this will not work as the link has /{str} system
		fmt.Println("hello world")
	default:
		conn.Write([]byte(response404))
	}
}
