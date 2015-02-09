package wol

import (
    "regexp"
    "errors"
    "strings"
    "encoding/hex"
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

// Returns a MacAddress object given a mac address string
func GetMacBytes(mac string) (*MacAddress, error) {
    // Parse MAC addr
    for _, delim := range delims {
        mac = strings.Replace(mac, string(delim), "", -1)
    }

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

// Constructs a "Magic Packet" broadcast frame which contains 6 bytes of
// 0xff followed by 16 repetitions of a given mac address.
//
// This function accepts a mac address string, and returns a pointer to
// a MagicPacket object */
func NewMagicPacket(mac string) (*MagicPacket, error) {
    var packet  MagicPacket

    // Parse the mac address into a "MacAddress". For the time being, only
    // the traditional methods of writing MAC addresses are supported.
    // XX:XX:XX:XX:XX or XX-XX-XX-XX-XX-XX will match. All other will throw
    // up an error to the caller.
    if re_MAC.MatchString(mac) {
        // Setup the header which is 6 repetitions of 0xff
        for idx, _ := range packet.header {
            packet.header[idx] = 0xff
        }

        addr, err := GetMacBytes(mac)
        if err != nil {
            return nil, err
        }

        // Setup the payload which is 6 repetitions of the MAC addr
        for idx, _ := range packet.payload {
            packet.payload[idx] = *addr
        }

        return &packet, nil
    }

    return nil, errors.New("Invalid MAC address format seen with " + mac)
}
