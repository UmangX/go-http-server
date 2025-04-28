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
	fmt.Printf("\n")
	fmt.Printf("-----------------------------------------\n")
	fmt.Printf("\n")

	line, _ := net.Listen("tcp", ":4221")
	for {
		conn, _ := line.Accept()
		go betterhandle(conn)
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

func betterhandle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	conn.Read(buf)

	// so the buff is the request done by the client
	buffer_content := string(buf)

	//content lines is the content seperated through \n
	content_lines := strings.Split(buffer_content, "\n")

	//for _, val := range content_lines {
	//fmt.Println(val)
	//}

	// there is the things which is need here
	// request type / url / header_data which is provided
	request_info := strings.Split(content_lines[0], " ")
	request_type := request_info[0]
	request_url_seperated := strings.Split(request_info[1], "/")
	request_url := request_url_seperated[1]

	if request_type == "GET" && request_url == "user-agent" {
		user_agent := ""
		for _, line := range content_lines {
			if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					user_agent = strings.TrimSpace(parts[1])
				}
				break
			}
		}
		writeResponse(conn, 200, user_agent)
		return
	}

	if request_type == "GET" && request_url == "echo" {

		if len(request_url_seperated) != 3 {
			writeResponse(conn, 200, " ")
			return
		}
		writeResponse(conn, 200, request_url_seperated[2])
		return
	}

	if request_url == "/" {
		writeResponse(conn, 200, "hello")
		return
	}

	writeResponse(conn, 404, " ")
}
