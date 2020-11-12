package wol

////////////////////////////////////////////////////////////////////////////////

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////////////////////////////

func TestNewMagicPacket(t *testing.T) {
	for _, tc := range []struct {
		mac      string
		expected MACAddress
	}{
		{"00:00:00:00:00:00", MACAddress{0, 0, 0, 0, 0, 0}},
		{"00:ff:01:03:00:00", MACAddress{0, 255, 1, 3, 0, 0}},
		{"00-ff-01-03-00-00", MACAddress{0, 255, 1, 3, 0, 0}},
	} {
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
	for _, tc := range []struct {
		mac string
	}{
		{"00x00:00:00:00:00"},
		{"00:00:Z0:00:00:00"},
		{"01:23:45:67:89:ab:cd:ef"},
		{"01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00"},
		{"01-23-45-67-89-ab-cd-ef"},
		{"01-23-45-67-89-ab-cd-ef-00-00-01-23-45-67-89-ab-cd-ef-00-00"},
		{"0123.4567.89ab"},
		{"0123.4567.89ab.cdef"},
		{"0123.4567.89ab.cdef.0000.0123.4567.89ab.cdef.0000"},
	} {
		_, err := New(tc.mac)
		assert.NotNil(t, err)
	}
}

func TestMagicPacketMarshal(t *testing.T) {
	for _, tc := range []struct {
		mac   string
		count int
	}{
		{"00:00:00:00:00:00", 102},
		{"00:ff:01:03:00:00", 102},
		{"00-ff-01-03-00-00", 102},
	} {
		pkt, err := New(tc.mac)
		assert.Equal(t, err, nil)

		bs, err := pkt.Marshal()
		assert.Equal(t, err, nil)

		assert.Equal(t, len(bs), tc.count)
	}
}
