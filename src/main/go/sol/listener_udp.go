package main

import (
	"net"
	"os"
	"strconv"
	"strings"
)

type MagicPacket []byte

func ListenerUDP(port int) {
	logger(3, "Listening UDP packets on port [" + strconv.Itoa(port) + "]")
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		logger(1, "Error while resolving local address: " + err.Error())
		os.Exit(1)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger(1, "Error while start listening: " + err.Error())
		os.Exit(1)
	}
	ReadPacket(sock)
}

func ReadPacket(sock *net.UDPConn) {
	var buf [1024]byte
	for {
		rlen, remote, err := sock.ReadFromUDP(buf[:])
		if err == nil {
			extractedMacAddress, _ := extractMacAddress(rlen, buf)
			logger(3, "Received a MAC address from IP [" + remote.String() + "], extracted mac [" + extractedMacAddress.String() + "]")
			if matchAddress(extractedMacAddress) {
				doAction()
			}
		} else {
			logger(2, "Error while reading a packet: " + err.Error())
		}
	}
}

func matchAddress(receivedAddress net.HardwareAddr) bool {
	receivedAddressAsString := receivedAddress.String()
	for _, value := range LocalNetworkMap() {
		if strings.HasPrefix(value, receivedAddressAsString) {
			return true
		}
	}

	return false
}

func extractMacAddress(rlen int, buf [1024]byte) (net.HardwareAddr, error) {
	var r = ""
	// TODO check whole magic packet structure (FF FF FF FF FF FF <MAC>*6)
	if rlen >= 12 {
		var sep = ""
		for i := 6; i < 12; i++ {
			val := int64(buf[i])                 // decimal value
			s := strconv.FormatInt(val, 16)      // convert to hexa (base 16)
			r = leftPad2Len(s, "0", 2) + sep + r // pad on two characters because some wake on lan tools are actually sending ":01:" as ":1:"
			sep = ":"
		}
	} else {
		logger(2, "The received packet is too small, size [" + strconv.Itoa(rlen) + "]")
	}
	return net.ParseMAC(r)
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func doAction() {
	for idx, _ := range configuration.Commands {
		Command := configuration.Commands[idx]
		if Command.Operation == configuration.Default {
			ExecuteCommand(Command)
			break
		}
	}
}
