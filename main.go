package main

import (
	"io"
	"log"
	"net"

	"github.com/nithishravindra/byo-lb-l4/internal/strategy"
)

func handleConnection(clientConn net.Conn, backendAddress string) {
	defer clientConn.Close()

	backendConn, err := net.Dial("tcp", backendAddress)

	if err != nil {
		log.Println("Error connecting to backend server:", err)
		return
	}
	defer backendConn.Close()
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

func getTcpConnection(req string) net.Listener {
	listener, err := net.Listen("tcp", req)
	if err != nil {
		log.Panic("Error creating TCP server:", err)
	}
	// Do not close the listener here
	return listener
}

func main() {
	log.Println("Initialized load balancer")
	backendAddressList := []string{"localhost:8080", "localhost:8081", "localhost:8082", "localhost:8083"}

	address := "localhost:9090"
	listener := getTcpConnection(address)
	defer listener.Close()
	rr := &strategy.RoundRobin{}
	// Accept incoming connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		backendAddress := rr.RoundRobin(backendAddressList)
		log.Printf("Listening on address: %v, forwarding to backend %v\n", address, backendAddress)
		// Handle the connection in a separate goroutine
		go handleConnection(clientConn, backendAddress)
	}
}
