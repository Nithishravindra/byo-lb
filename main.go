package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/nithishravindra/byo-lb-l4/internal/ratelimiter"
	"github.com/nithishravindra/byo-lb-l4/internal/strategy"
)

type JsonErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func handleConnection(clientConn net.Conn, backendAddress string, allowed bool) {
	defer clientConn.Close()

	backendConn, err := net.Dial("tcp", backendAddress)

	// If not allowed, send a message and return
	if !allowed {
		response := JsonErrorResponse{Status: 429, Error: "Request limit exceeded. Please try again later."}
		responseJson, _ := json.Marshal(response)

		clientConn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n"))
		clientConn.Write([]byte("Content-Type: application/json\r\n"))
		clientConn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(responseJson))))
		clientConn.Write([]byte("\r\n"))
		clientConn.Write(responseJson)
		return
	}

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

	var rl ratelimiter.RateLimiter

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		backendAddress := rr.RoundRobin(backendAddressList)
		log.Printf("Listening on address: %v, forwarding to backend %v\n", address, backendAddress)
		go handleConnection(clientConn, backendAddress, rl.Allow())

	}
}
