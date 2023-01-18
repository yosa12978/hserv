package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleIncomingRequest(conn)
	}
	// req := ParseRequest(`GET / HTTP/1.1
	// Host: localhost:8089
	// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0
	// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8

	// <h1>Hello World</h1>
	// `)
	// fmt.Printf("%s\n", req.Version)
	// fmt.Printf("%s\n", req.Method)
	// fmt.Printf("%s\n", req.Path)
	// fmt.Printf("%s\n", req.Body)
}

func handleIncomingRequest(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1<<10)
	_, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		log.Printf(err.Error())
	}
	filename := "index.html"
	file, err := os.ReadFile("./www/" + filename)
	if err != nil {
		e404 := "<h1>error 404: file not found</h1>"
		fmt.Fprintf(conn, "HTTP/1.1 404 NOT FOUND\r\n")
		fmt.Fprintf(conn, "Content-Length: %d\r\n", len(e404))
		fmt.Fprintf(conn, "Content-Type: text/html\r\n")
		fmt.Fprintf(conn, "\r\n")
		fmt.Fprintf(conn, e404)
		return
	}
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(file))
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, string(file))
	fmt.Println(string(buffer))
}

type Header struct {
	Key   string
	Value string
}

func (h *Header) Format() string {
	return h.Key + ": " + h.Value
}

type Response struct {
	Version string
	Status  Status
	Headers []Header
	Body    string
}

func (res *Response) ToString() string {
	return ""
}

type Status struct {
	Code    int
	Message string
}

type Request struct {
	Version string
	Method  string
	Path    string
	Headers []Header
	Body    string
}

func ParseRequest(s string) *Request {
	lines := strings.Split(s, "\n")
	fl := strings.Split(lines[0], " ")
	req := &Request{}
	req.Method = fl[0]
	req.Path = fl[1]
	req.Version = fl[2]
	body := strings.Split(s, "\r\n")[1:]
	req.Body = strings.Join(body, "\r\n")
	return req
}
