package main

import (
	"log"

	"github.com/linde12/gowol"
)

func startMachineX(mac string) {
	err := startMachine(mac)
	if err != nil {
		log.Printf("Failed to start %s: %v", mac, err)
	}
}

func startMachine(mac string) error {
	packet, err := gowol.NewMagicPacket(mac)
	if err != nil {
		return err
	}
	// TODO: add selecting port and IP
	log.Printf("Starting machine: %s", mac)
	return packet.Send("255.255.255.255")
}
