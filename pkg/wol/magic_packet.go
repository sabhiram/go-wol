package wol

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
)

////////////////////////////////////////////////////////////////////////////////

var (
	delims = ":-"
	reMAC  = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delims + `]){5}([0-9a-fA-F]{2})$`)
)

////////////////////////////////////////////////////////////////////////////////

// MACAddress6Byte represents a 6 byte network mac address.
type MACAddress6Byte [6]byte

type Packet struct {
	header  [6]byte
	payload [16]MACAddress6Byte
}

// MagicPacket is constituted of 6 bytes of 0xFF followed by 16-groups of the
// destination MAC address.
// Implements `encoding.BinaryMarshaler`
// See: https://en.wikipedia.org/wiki/Wake-on-LAN#Magic_packet
type MagicPacket struct {
	MacAddr   string
	BcastAddr string

	Packet Packet
}

// NewMagicPacket returns a magic packet based on a mac address string.
func NewMagicPacket(macAddr string) (*MagicPacket, error) {
	macAddr6Bytes := MACAddress6Byte{}

	mp := MagicPacket{
		MacAddr: macAddr,
		Packet:  Packet{},
	}

	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return nil, err
	}

	// We only support 6 byte MAC addresses since it is much harder to use the
	// binary.Write(...) interface when the size of the MagicPacket is dynamic.
	if !reMAC.MatchString(macAddr) {
		return nil, fmt.Errorf("%s is not a IEEE 802 MAC-48 address", macAddr)
	}

	// Copy bytes from the returned HardwareAddr -> a fixed size MACAddress.
	for idx := range macAddr6Bytes {
		macAddr6Bytes[idx] = hwAddr[idx]
	}

	// Setup the header which is 6 repetitions of 0xFF.
	for idx := range mp.Packet.header {
		mp.Packet.header[idx] = 0xFF
	}

	// Setup the payload which is 16 repetitions of the MAC addr.
	for idx := range mp.Packet.payload {
		mp.Packet.payload[idx] = macAddr6Bytes
	}

	return &mp, nil
}

// MarshalBinary serializes the magic packet structure into a 102 byte slice.
func (mp *MagicPacket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// by isolating the real payload in `.Packet`
	// we can add extra fields to the struct
	if err := binary.Write(&buf, binary.BigEndian, mp.Packet); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
