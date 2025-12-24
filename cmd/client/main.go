package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/protocol"
	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/tunnel"
)

var (
	serverAddr = flag.String("server", "192.168.88.185:6969", "Server address")
	clientID   = flag.String("id", "test-client", "Client identifier")
)

func main() {
	flag.Parse()
	log.Printf("Starting VPN Client (ID: %s)", *clientID)
	log.Printf("Connecting to server: %s", *serverAddr)

	// Connect to server
	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server")

	// Set timeouts
	conn.SetDeadline(time.Time{}) // No timeout for now

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Send hello
	helloPacket := tunnel.CreateHelloPacket(*clientID)
	if err := tunnel.WritePacket(conn, helloPacket); err != nil {
		log.Fatalf("Failed to send hello: %v", err)
	}
	log.Println("Sent hello to server")

	// Wait group for goroutines
	var wg sync.WaitGroup

	// Start reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		readFromServer(conn)
	}()

	// Start writer goroutine (send periodic pings)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sendPeriodicPings(conn, sigCh)
	}()

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down client...")

	// Close connection to stop goroutines
	conn.Close()

	// Wait for goroutines to finish
	wg.Wait()
	log.Println("Client stopped")
}

func readFromServer(conn net.Conn) {
	packetCount := 0

	for {
		packet, err := tunnel.ReadPacket(conn)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Read error from server: %v", err)
			}
			return
		}

		packetCount++
		processServerPacket(packet, packetCount)
	}
}

func processServerPacket(packet *protocol.Packet, count int) {
	ts := time.Unix(int64(packet.Header.Timestamp), 0)
	latency := time.Since(ts)

	log.Printf("\n=== Packet #%d from server ===", count)
	log.Printf("Type:      %d", packet.Header.Type)
	log.Printf("Seq:       %d", packet.Header.Seq)
	log.Printf("Latency:   %v", latency)

	switch packet.Header.Type {
	case protocol.TypeData:
		log.Printf("Server says: %s", string(packet.Payload))

	case protocol.TypePong:
		log.Printf("Received PONG for seq %d", packet.Header.Seq)
		log.Printf("Payload: %s", string(packet.Payload))

	default:
		log.Printf("Unknown type %d, payload: %s",
			packet.Header.Type, string(packet.Payload))
	}

	log.Printf("===========================\n")
}

func sendPeriodicPings(conn net.Conn, stopCh chan os.Signal) {
	seq := uint32(0)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			seq++
			pingPacket := tunnel.CreatePingPacket(seq)

			log.Printf("Sending PING #%d to server", seq)

			if err := tunnel.WritePacket(conn, pingPacket); err != nil {
				log.Printf("Failed to send PING: %v", err)
				return
			}

		case <-stopCh:
			return
		}
	}
}
