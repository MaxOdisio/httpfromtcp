package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
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
			ch := getLinesChannel(c)
			for line := range ch {
				fmt.Printf("%s\n", line)
			}
			fmt.Printf("Connection closed for %s\n", c.RemoteAddr())
		}(conn)

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		buffer := make([]byte, 8)
		var currentLine []byte

		for {
			n, err := f.Read(buffer)
			if n > 0 {
				currentLine = append(currentLine, buffer[:n]...)

				// creao un ciclo per verificare se sono presenti più caratteri "\n" nel buffer
				for {
					i := bytes.Index(currentLine, []byte("\n"))

					if i == -1 {
						break
					}

					line := string(currentLine[:i])
					ch <- line
					currentLine = currentLine[i+1:]
				}
			}

			if err != nil {
				if err == io.EOF && len(currentLine) > 0 {
					ch <- string(currentLine)
				}
				break
			}
		}
	}()

	return ch
}
