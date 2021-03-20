package wol

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

// MacIface holds a MAC Address to wake up, along with an optionally specified
// default interface to use when typically waking up said interface.
type MacIface struct {
	Mac   string
	Iface string
}

// DecodeToMacIface takes a byte buffer and converts decodes it using the gob
// package to a MacIface entry.
func DecodeToMacIface(buf *bytes.Buffer) (MacIface, error) {
	var entry MacIface
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&entry)
	return entry, err
}

// EncodeFromMacIface takes a MAC and an Iface and encodes a gob with a MacIface
// entry.
func EncodeFromMacIface(mac, iface string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	entry := MacIface{mac, iface}
	err := gob.NewEncoder(buf).Encode(entry)
	return buf, err
}

// IPFromInterface returns a `*net.UDPAddr` from a network interface name.
func IPFromInterface(iface string) (*net.UDPAddr, error) {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}

	addrs, err := ief.Addrs()
	if err == nil && len(addrs) <= 0 {
		err = fmt.Errorf("no address associated with interface %s", iface)
	}
	if err != nil {
		return nil, err
	}

	// Validate that one of the addrs is a valid network IP address.
	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return &net.UDPAddr{
					IP: ip.IP,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no address associated with interface %s", iface)
}
