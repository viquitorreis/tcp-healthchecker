package main

import (
	"log"
	"net"
	"time"
)

func main() {
	TCPListener()
}

func TCPListener() {
	listener, err := net.Listen("tcp", "127.0.0.1:4321")
	if err != nil {
		log.Fatalf("net.Listen() error: %v", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("listener.Close() error: %v", err)
		}
	}()

	log.Printf("Listening on %s", listener.Addr())

	// handle new connections / accept new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept() error: %v", err)
			continue
		}

		// handle / process the connection accepted
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	now := time.Now()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn.Close() error: %v", err)
		}
		log.Printf("Close connection from %s. Connection duration: %v ms", conn.RemoteAddr(), time.Since(now))
	}()

	// connection accepted - processing it
	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	_, err := conn.Write([]byte("Hello, client!"))
	if err != nil {
		log.Printf("conn.Write() error: %v", err)
	}
}
