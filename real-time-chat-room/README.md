# TCP Chat Server in Go

## Introduction

This document serves as a comprehensive guide for building a simple TCP chat server in Go. The server allows multiple clients to connect, send messages, and receive messages in real-time.

## Key Concepts

### 1. TCP Protocol

- **Transmission Control Protocol (TCP)** is a connection-oriented protocol used for reliable communication between client and server.
- TCP ensures that data is sent and received in the correct order and without duplication.

### 2. Go's `net` Package

- The `net` package in Go provides a portable interface for network programming.
- It includes functionalities to create servers and clients, handle connections, and manage network I/O.

### 3. Goroutines and Concurrency

- **Goroutines** are lightweight threads managed by the Go runtime.
- The server uses goroutines to handle multiple client connections concurrently, allowing each client to communicate with the server without blocking others.

### 4. Broadcasting Messages

- When a client sends a message, the server broadcasts it to all connected clients.
- This involves iterating over the list of connected clients and sending the message to each one.

### 5. Mutex for Synchronization

- A **mutex** (mutual exclusion) is used to protect shared resources (like the list of connections) from concurrent access issues.
- This ensures that only one goroutine can modify the list of clients at a time, preventing race conditions.

## Implementation Steps

### Step 1: Setting Up the TCP Server

1. **Listen for Connections**:
   Use `net.Listen` to create a TCP server that listens on a specific port.

   ```go
   listener, err := net.Listen("tcp", ":8080")
   ```

2. **Accept Incoming Connections**:
   Use a loop to continuously accept new client connections.

   ```go
   for {
       conn, err := listener.Accept()
       if err != nil {
           log.Println("Error accepting connection:", err)
           continue
       }
       go handleConnection(conn) // Handle connection in a separate goroutine
   }
   ```

### Step 2: Handling Client Connections

1. **Define `handleConnection` Function**:
   This function will be responsible for reading messages from the client and broadcasting them.

   ```go
   func handleConnection(conn net.Conn) {
       defer conn.Close() // Ensure the connection is closed when done
       // (Add logic to read messages and broadcast them)
   }
   ```

2. **Reading Messages**:
   Use a buffer to read messages from the client. You may want to handle potential errors while reading.

   ```go
   buffer := make([]byte, 1024)
   n, err := conn.Read(buffer)
   ```

### Step 3: Broadcasting Messages

1. **Store Connected Clients**:
   Maintain a list of active client connections to broadcast messages.

   ```go
   var connections []net.Conn
   ```

2. **Broadcast Function**:
   Create a function to send a message to all connected clients.

   ```go
   func broadcast(message string) {
       for _, client := range connections {
           _, err := client.Write([]byte(message))
           if err != nil {
               log.Println("Error sending message:", err)
           }
       }
   }
   ```

### Step 4: Synchronizing Access to Shared Resources

1. **Using Mutex**:
   Protect the list of connections with a mutex to prevent concurrent modifications.

   ```go
   var connMutex sync.Mutex
   ```

   Use `connMutex.Lock()` and `connMutex.Unlock()` when modifying the connections list.

### Example Code Snippet

Here is a simplified version of how the core components fit together:

```go
func main() {
   listener, err := net.Listen("tcp", ":8080")
   if err != nil {
       log.Fatal(err)
   }
   defer listener.Close()

   for {
       conn, err := listener.Accept()
       if err != nil {
           log.Println("Error accepting connection:", err)
           continue
       }
       go handleConnection(conn)
   }
}
```

### Additional Considerations

- **Error Handling**: Implement robust error handling throughout the server to manage connection issues and message transmission failures.
- **Graceful Shutdown**: Consider how to cleanly shut down the server and close all client connections when needed.

## Conclusion

By following these steps and understanding the key concepts, you can build a basic TCP chat server in Go that can handle multiple clients concurrently. This server serves as a foundation for further enhancements, such as adding private messaging, improved error handling, or user authentication in the future.
