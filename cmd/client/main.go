package main

import (
	"log"

	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/cmd"
	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/crypto"
	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/transport"
	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/tun"
)

func main() {
	tunDev := tun.Create()
	config := cmd.LoadConfig()

	udp, err := transport.Dial(config.ServerAddress + config.ListenPort)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, config.MTU)

	for {
		n, err := tunDev.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d bytes from TUN\n", n)

		encrypted, err := crypto.Encrypt(config.EncryptionKey, buf[:n])
		if err != nil {
			continue
		}

		udp.Conn.Write(encrypted)
	}
}
