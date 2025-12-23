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

	udp, err := transport.Listen("0.0.0.0" + config.ListenPort)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2048)

	for {
		n, addr, err := udp.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		log.Printf("Received %d encrypted bytes\n", n)

		packet, err := crypto.Decrypt(config.EncryptionKey, buf[:n])
		if err != nil {
			continue
		}

		tunDev.Write(packet)

		reply := make([]byte, config.MTU)
		rn, err := tunDev.Read(reply)
		if err == nil {
			enc, _ := crypto.Encrypt(config.EncryptionKey, reply[:rn])
			udp.Conn.WriteTo(enc, addr)
		}
	}
}
