package main

import (
	"fmt"
	"net"

	"github.com/maxodisio/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error on tcp listener: %v\n", err)
		return
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Printf("Error creating tcp connection: %v\n", err)
			continue
		}

		go func(c net.Conn) {
			fmt.Printf("Connection established from %s\n", c.RemoteAddr())
			req, err := request.RequestFromReader(c)
			if err != nil {
				fmt.Printf("Error on request from reader: %v\n", err)
				return
			}
			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
			if len(req.Headers) != 0 {
				fmt.Println("Headers:")
				for k, v := range req.Headers {
					fmt.Printf("- %s: %s\n", k, v)
				}
			}
			fmt.Println("===========================")
			fmt.Printf("Connection closed for %s\n", c.RemoteAddr())
		}(conn)

	}

}
