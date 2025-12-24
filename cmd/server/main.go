package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/protocol"
	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/tunnel"
)

const (
	serverPort = ":6969"
)

func main() {
	log.Println("Starting VPN Server...")
	log.Printf("Listening on port %s", serverPort)

	// Start TCP server
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				continue
			}
			go handleClient(conn)
		}
	}()

	log.Println("Server is ready. Press Ctrl+C to stop.")
	<-sigCh
	log.Println("Shutting down server...")
}

func handleClient(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	log.Printf("New client connected: %s", clientAddr)
	defer func() {
		conn.Close()
		log.Printf("Client disconnected: %s", clientAddr)
	}()

	// Set timeouts
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// Handle client packets
	packetCount := 0

	for {
		// Read packet
		packet, err := tunnel.ReadPacket(conn)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Read error from %s: %v", clientAddr, err)
			}
			return
		}

		packetCount++
		processPacket(conn, packet, clientAddr, packetCount)
	}
}

func processPacket(conn net.Conn, packet *protocol.Packet, clientAddr string, count int) {
	// Log packet info
	ts := time.Unix(int64(packet.Header.Timestamp), 0)
	latency := time.Since(ts)

	log.Printf("\n=== Packet #%d from %s ===", count, clientAddr)
	log.Printf("Type:      %d", packet.Header.Type)
	log.Printf("Seq:       %d", packet.Header.Seq)
	log.Printf("Length:    %d bytes", packet.Header.Length)
	log.Printf("Timestamp: %s", ts.Format("15:04:05"))
	log.Printf("Latency:   %v", latency)
	log.Printf("Payload:   %s", string(packet.Payload))

	// Handle different packet types
	switch packet.Header.Type {
	case protocol.TypeHello:
		log.Printf("Client says: %s", string(packet.Payload))
		// Send welcome message
		welcome := &protocol.Packet{
			Header: protocol.Header{
				Type:      protocol.TypeData,
				Length:    16,
				Seq:       packet.Header.Seq,
				Timestamp: uint32(time.Now().Unix()),
			},
			Payload: []byte("Welcome to VPN!"),
		}
		if err := tunnel.WritePacket(conn, welcome); err != nil {
			log.Printf("Failed to send welcome: %v", err)
		}

	case protocol.TypePing:
		log.Printf("Received PING from %s", clientAddr)
		// Send pong response
		pong := tunnel.CreatePongPacket(packet)
		if err := tunnel.WritePacket(conn, pong); err != nil {
			log.Printf("Failed to send PONG: %v", err)
		}

	case protocol.TypeData:
		log.Printf("Data packet from %s", clientAddr)
		// Echo back
		echo := &protocol.Packet{
			Header: protocol.Header{
				Type:      protocol.TypeData,
				Length:    packet.Header.Length,
				Seq:       packet.Header.Seq,
				Timestamp: uint32(time.Now().Unix()),
			},
			Payload: append([]byte("Echo: "), packet.Payload...),
		}
		if err := tunnel.WritePacket(conn, echo); err != nil {
			log.Printf("Failed to echo: %v", err)
		}

	default:
		log.Printf("Unknown packet type: %d", packet.Header.Type)
	}

	log.Printf("=====================\n")
}
