package wol

import (
	"github.com/stretchr/testify/assert"
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
		assert.NotEqual(test, err, nil)
	}
}
