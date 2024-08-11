package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func handleConnection(clientConn net.Conn, backendAddress string) {
	defer clientConn.Close()

	// Dial the backend server
	backendConn, err := net.Dial("tcp", backendAddress)
	if err != nil {
		log.Println("Error connecting to backend server:", err)
		return
	}
	defer backendConn.Close()

	fmt.Println("Connected to backend server at", backendAddress)

	// Channel to signal when either direction is done
	done := make(chan struct{})

	// Forward data from client to backend
	go func() {
		io.Copy(backendConn, clientConn)
		done <- struct{}{}
	}()

	// Forward data from backend to client
	go func() {
		io.Copy(clientConn, backendConn)
		done <- struct{}{}
	}()

	// Wait for either direction to finish
	<-done
}

func main() {
	log.Println("Initialized load balancer")
	address := "localhost:9090"
	backendAddress := "localhost:8080"

	// Start listening on the specified address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic("Error creating TCP server:", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on address: %v, forwarding to backend: %v\n", address, backendAddress)

	// Accept incoming connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(clientConn, backendAddress)
	}
}
