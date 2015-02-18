package wol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Define globals for the MagicPacket and the MacAddress parsing
var (
	delims = ":-"
	re_MAC = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delims + `]){5}([0-9a-fA-F]{2})$`)
)

// A MacAddress is 6 bytes in a row
type MacAddress [6]byte

// A MagicPacket is constituted of 6 bytes of 0xFF followed by
// 16 groups of the destination MAC address.
type MagicPacket struct {
	header  [6]byte
	payload [16]MacAddress
}

// Returns a pointer to a MacAddress, given a valid MAC Address string
func getMacAddressFromString(mac string) (*MacAddress, error) {
	// First strip the delimiters from the valid MAC Address
	for _, delim := range delims {
		mac = strings.Replace(mac, string(delim), "", -1)
	}

	// Fetch the bytes from the string representation of the
	// MAC address. Address is []byte
	address, err := hex.DecodeString(mac)
	if err != nil {
		return nil, err
	}

	var ret MacAddress
	for idx := range ret {
		ret[idx] = address[idx]
	}

	return &ret, nil
}

// This function accepts a MAC Address string, and returns a pointer to
// a MagicPacket object. A Magic Packet is a broadcast frame which
// contains 6 bytes of 0xFF followed by 16 repetitions of a given mac address.
func NewMagicPacket(mac string) (*MagicPacket, error) {
	var packet MagicPacket

	// Parse the MAC Address into a "MacAddress". For the time being, only
	// the traditional methods of writing MAC Addresses are supported.
	// XX:XX:XX:XX:XX or XX-XX-XX-XX-XX-XX will match. All others will throw
	// up an error to the caller.
	if re_MAC.MatchString(mac) {
		// Setup the header which is 6 repetitions of 0xFF
		for idx := range packet.header {
			packet.header[idx] = 0xFF
		}

		addr, err := getMacAddressFromString(mac)
		if err != nil {
			return nil, err
		}

		// Setup the payload which is 16 repetitions of the MAC addr
		for idx := range packet.payload {
			packet.payload[idx] = *addr
		}

		return &packet, nil
	}
	return nil, errors.New("Invalid MAC address format seen with " + mac)
}

// This function accepts a MAC address string, and s
// Function to send a magic packet to a given mac address
func SendMagicPacket(macAddr, bcastAddr string) error {
	magicPacket, err := NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, magicPacket)

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		fmt.Printf("Unable to get a UDP address for %s\n", bcastAddr)
		return err
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Unable to dial UDP address for %s\n", bcastAddr)
		return err
	}
	defer connection.Close()

	bytesWritten, err := connection.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Unable to write packet to connection\n")
		return err
	} else if bytesWritten != 102 {
		fmt.Printf("Warning: %d bytes written, %d expected!\n", bytesWritten, 102)
	}

	return nil
}
