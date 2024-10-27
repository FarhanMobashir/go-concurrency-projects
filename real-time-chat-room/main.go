package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func handleConnection(conn net.Conn) {
	defer func() {
		// remove connection from the list when done
		connMutex.Lock()
		for i, c := range connections {
			if c == conn {
				connections = append(connections[:i], connections[i+1:]...)
			}
		}
		connMutex.Unlock()
		conn.Close()
	}()

	fmt.Println("New Client Connected:", conn.RemoteAddr())

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client Disconnected:", conn.RemoteAddr())
			return
		}

		// Construct the message string
		message := fmt.Sprintf("%s: %s", conn.RemoteAddr(), string(buffer[:n]))

		// Broadcast the message to all clients
		broadcast(message)
		fmt.Printf("Received message from %s: %s\n", conn.RemoteAddr(), message)
	}
}

// broadcast function
func broadcast(message string) {
	connMutex.Lock()
	defer connMutex.Unlock()

	for _, conn := range connections {
		_, err := conn.Write([]byte(message))

		if err != nil {
			fmt.Println("Error broadcasting the messages to the client", err)
		}
	}
}

var (
	connections []net.Conn // slice to hold all the active connection
	connMutex   sync.Mutex // Mutex to manage concurrent access to all connection
)

func main() {
	fmt.Println("Hello from chat room")
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		// add the connection to the connection list
		connMutex.Lock()
		connections = append(connections, conn)
		connMutex.Unlock()

		// go routine
		go handleConnection(conn)
	}
}
