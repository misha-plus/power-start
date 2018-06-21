package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/linde12/gowol"
)

func startMachine(mac string) error {
	packet, err := gowol.NewMagicPacket(mac)
	if err != nil {
		return err
	}

	log.Printf("Starting machine: MAC = %s", mac)
	_, network, err := net.ParseCIDR(config.Network)
	if err != nil {
		return err
	}

	if len(network.IP) != 4 {
		return errors.New("Only IPv4 is supported")
	}

	toUint32 := func(buf []byte) uint32 {
		var y uint32
		_ = binary.Read(bytes.NewReader(buf), binary.BigEndian, &y)
		return y
	}

	broadcast := make([]byte, 4)
	binary.BigEndian.PutUint32(broadcast, toUint32(network.IP)|^toUint32(network.Mask))
	return packet.Send(net.IP(broadcast).String())
}
