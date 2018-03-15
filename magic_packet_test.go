package wol

////////////////////////////////////////////////////////////////////////////////

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////////////////////////////

func TestNewMagicPacket(t *testing.T) {
	var PositiveTestCases = []struct {
		mac      string
		expected MACAddress
	}{
		{"00:00:00:00:00:00", MACAddress{0, 0, 0, 0, 0, 0}},
		{"00:ff:01:03:00:00", MACAddress{0, 255, 1, 3, 0, 0}},
		{"00-ff-01-03-00-00", MACAddress{0, 255, 1, 3, 0, 0}},
	}

	for _, tc := range PositiveTestCases {
		pkt, err := New(tc.mac)
		for _, v := range pkt.header {
			assert.Equal(t, int(v), 255)
		}
		for _, mac := range pkt.payload {
			assert.Equal(t, tc.expected, mac)
		}
		assert.Equal(t, err, nil)
	}
}

func TestNewMagicPacketNegative(t *testing.T) {
	var NegativeTestCases = []struct {
		mac string
	}{
		{"00x00:00:00:00:00"},
		{"00:00:Z0:00:00:00"},
	}

	for _, tc := range NegativeTestCases {
		_, err := New(tc.mac)
		assert.NotNil(t, err)
	}
}

func TestGetIpFromInterface(t *testing.T) {
	interfaces, err := net.Interfaces()
	assert.Nil(t, err)

	// We can't actually enforce that we get a valid IP, but
	// either the error or the pointer should be nil
	for _, i := range interfaces {
		addr, err := GetIpFromInterface(i.Name)
		if err == nil {
			assert.NotNil(t, addr)
		} else {
			assert.Nil(t, addr)
		}
	}
}

func TestGetIpFromInterfaceNegative(t *testing.T) {
	// Test some fake interfaces
	var NegativeTestCases = []struct {
		iface string
	}{
		{"fake-interface-0"},
		{"fake-interface-1"},
	}

	for _, tc := range NegativeTestCases {
		addr, err := GetIpFromInterface(tc.iface)
		assert.Nil(t, addr)
		assert.NotNil(t, err)
	}
}

////////////////////////////////////////////////////////////////////////////////
