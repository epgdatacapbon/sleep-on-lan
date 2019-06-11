package main

import (
	"encoding/hex"
	"errors"
	"net"
	"strings"
)

// Use the net library to return all Interfaces
// and capture any errors.
func GetInterfaces() []net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Warning("Unable to get interfaces: " + err.Error())
	}
	return interfaces
}

func LocalNetworkMap() map[string]string {
	result := make(map[string]string)
	for _, inter := range GetInterfaces() {
		addresses, _ := inter.Addrs()
		for _, addr := range addresses {
			result[addr.String()] = inter.HardwareAddr.String()
		}
	}
	return result
}

// Use a MAC address to form a magic packet
// macAddr form 12:34:56:78:9a:bc
func EncodeMagicPacket(macAddr string) (MagicPacket, error) {
	if len(macAddr) != (6*2 + 5) {
		return nil, errors.New("Invalid MAC Address [" + macAddr + "]")
	}

	macBytes, err := hex.DecodeString(strings.Join(strings.Split(macAddr, ":"), ""))
	if err != nil {
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}

	return MagicPacket(b), nil
}

// Send a Magic Packet to an broadcast class IP address via UDP
func (p MagicPacket) Wake(bcastAddr string) error {
	a, err := net.ResolveUDPAddr("udp", bcastAddr+":9")
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}

	written, err := c.Write(p)
	c.Close()

	// Packet must be 102 bytes in length
	if written != 102 {
		return err
	}

	return nil
}
