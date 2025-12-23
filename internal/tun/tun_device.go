package tun

import (
	"log"

	"github.com/songgao/water"
)

func Create() *water.Interface {
	cfg := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(cfg)
	if err != nil {
		log.Fatalf("failed to create TUN: %v", err)
	}

	log.Printf("TUN device created: %s\n", iface.Name())
	return iface
}
