package wol

import (
	"encoding/hex"
	"errors"
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
func GetMacAddressFromString(mac string) (*MacAddress, error) {
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
	for idx, _ := range ret {
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
		for idx, _ := range packet.header {
			packet.header[idx] = 0xFF
		}

		addr, err := GetMacAddressFromString(mac)
		if err != nil {
			return nil, err
		}

		// Setup the payload which is 16 repetitions of the MAC addr
		for idx, _ := range packet.payload {
			packet.payload[idx] = *addr
		}

		return &packet, nil
	}
	return nil, errors.New("Invalid MAC address format seen with " + mac)
}
