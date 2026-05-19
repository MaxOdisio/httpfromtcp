package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const address = "localhost:42069"

func main() {
	remoteAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, remoteAddress)
	if err != nil {
		fmt.Println("Error dialing UDP connection:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to", address, "- Write something and press Enter")

	for {
		fmt.Print("> ")

		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nAll read (EOF)")
				break
			}
			fmt.Println("Error reading input:", err)
			break
		}

		cleanS := strings.TrimSpace(s) // rimuoviamo il carattere "\n"

		// se l'utente preme solo invio, saltiamo
		if cleanS == "" {
			continue
		}

		_, err = conn.Write([]byte(cleanS))
		if err != nil {
			fmt.Printf("Error writing line <<%s>>: %v\n", s, err)
			break
		}
	}
}
