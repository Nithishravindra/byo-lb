package main

package main

import (
	"fmt"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}
		fmt.Println("Received:", string(buffer[:n]))
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			break
		}
	}
}

func createTCPServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error creating TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("TCP server listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	address := "localhost:8080"
	createTCPServer(address)
}