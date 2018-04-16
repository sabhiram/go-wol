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

// MACAddress represents a 6 byte network mac address.
type MACAddress [6]byte

// A MagicPacket is constituted of 6 bytes of 0xFF followed by 16-groups of the
// destination MAC address.
type MagicPacket struct {
	header  [6]byte
	payload [16]MACAddress
}

// New returns a magic packet based on a mac address string.
func New(mac string) (*MagicPacket, error) {
	var packet MagicPacket
	var macAddr MACAddress

	// We only support 6 byte MAC addresses since it is much harder to use the
	// binary.Write(...) interface when the size of the MagicPacket is dynamic.
	if !reMAC.MatchString(mac) {
		return nil, fmt.Errorf("invalid mac-address %s", mac)
	}

	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// Copy bytes from the returned HardwareAddr -> a fixed size MACAddress.
	for idx := range macAddr {
		macAddr[idx] = hwAddr[idx]
	}

	// Setup the header which is 6 repetitions of 0xFF.
	for idx := range packet.header {
		packet.header[idx] = 0xFF
	}

	// Setup the payload which is 16 repetitions of the MAC addr.
	for idx := range packet.payload {
		packet.payload[idx] = macAddr
	}

	return &packet, nil
}

////////////////////////////////////////////////////////////////////////////////

// GetIPFromInterface returns a `*net.UDPAddr` from a network interface name.
func GetIPFromInterface(iface string) (*net.UDPAddr, error) {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	addrs, err := ief.Addrs()
	if err != nil {
		return nil, err
	} else if len(addrs) <= 0 {
		return nil, fmt.Errorf("no address associated with interface %s", iface)
	}

	// Validate that one of the addr's is a valid network IP address.
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			// Verify that the DefaultMask for the address we want to use exists.
			if ip.IP.DefaultMask() != nil {
				return &net.UDPAddr{
					IP: ip.IP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no address associated with interface %s", iface)
}

// SendMagicPacket sends a magic packet UDP broadcast to the specified `macAddr`.
// The broadcast is sent to `bcastAddr` via the `iface`.  An empty `iface`
// implies a nil local address to dial.
func SendMagicPacket(macAddr, bcastAddr, iface string) error {
	// Construct a MagicPacket for the given MAC Address.
	magicPacket, err := New(macAddr)
	if err != nil {
		return err
	}

	// Fill our byte buffer with the bytes in our MagicPacket.
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, magicPacket)
	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)

	// Get a UDPAddr to send the broadcast to.
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// If an interface was specified, get the address associated with it.
	var localAddr *net.UDPAddr
	if iface != "" {
		var err error
		localAddr, err = GetIPFromInterface(iface)
		if err != nil {
			return err
		}
	}

	// Open a UDP connection, and defer it's cleanup.
	connection, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer connection.Close()

	// Write the bytes of the MagicPacket to the connection.
	n, err := connection.Write(buf.Bytes())
	if err != nil {
		return err
	} else if n != 102 {
		fmt.Printf("Warning: %d bytes written, %d expected!\n", n, 102)
	}

	return nil
}
