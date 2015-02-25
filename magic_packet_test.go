package wol

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestNewMagicPacket(test *testing.T) {
	var PositiveTestCases = []struct {
		mac      string
		expected MACAddress
	}{
		{"00:00:00:00:00:00", MACAddress{0, 0, 0, 0, 0, 0}},
		{"00:ff:01:03:00:00", MACAddress{0, 255, 1, 3, 0, 0}},
		{"00-ff-01-03-00-00", MACAddress{0, 255, 1, 3, 0, 0}},
	}

	for _, t := range PositiveTestCases {
		pkt, err := NewMagicPacket(t.mac)
		for _, v := range pkt.header {
			assert.Equal(test, int(v), 255)
		}
		for _, mac := range pkt.payload {
			assert.Equal(test, t.expected, mac)
		}
		assert.Equal(test, err, nil)
	}
}

func TestNewMagicPacketNegative(test *testing.T) {
	var NegativeTestCases = []struct {
		mac string
	}{
		{"00x00:00:00:00:00"},
		{"00:00:Z0:00:00:00"},
	}

	for _, t := range NegativeTestCases {
		_, err := NewMagicPacket(t.mac)
		assert.NotNil(test, err)
	}
}

func TestGetIpFromInterface(test *testing.T) {
	interfaces, err := net.Interfaces()
	assert.Nil(test, err)

	// We can't actually enforce that we get a valid IP, but
	// either the error or the pointer should be nil
	for _, i := range interfaces {
		addr, err := GetIpFromInterface(i.Name)
		if err == nil {
			assert.NotNil(test, addr)
		} else {
			assert.Nil(test, addr)
		}
	}
}

func TestGetIpFromInterfaceNegative(test *testing.T) {
	// Test some fake interfaces
	var NegativeTestCases = []struct {
		iface string
	}{
		{"fake-interface-0"},
		{"fake-interface-1"},
	}

	for _, t := range NegativeTestCases {
		addr, err := GetIpFromInterface(t.iface)
		assert.Nil(test, addr)
		assert.NotNil(test, err)
	}
}
