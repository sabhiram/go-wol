package wol

import (
	"bytes"
	"encoding/gob"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate the DecodeToMacIface function.
func TestDecodeToMacIface(t *testing.T) {
	var TestCases = []MacIface{
		{Mac: "00:00:00:00:00:00", Iface: ""},
		{Mac: "00:00:00:00:00:AA", Iface: "eth1"},
	}

	for _, entry := range TestCases {
		// First encode the MacIface to a bunch of bytes.
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(entry)
		assert.Nil(t, err)

		result, err := DecodeToMacIface(buf)
		assert.Nil(t, err)
		assert.Equal(t, entry.Mac, result.Mac)
		assert.Equal(t, entry.Iface, result.Iface)
	}
}

// Validate the EncodeFromMacIface function.
func TestEncodeFromMacIface(t *testing.T) {
	var TestCases = []MacIface{
		{Mac: "00:00:00:00:00:00", Iface: "eth0"},
		{Mac: "00:00:00:00:00:AA", Iface: ""},
	}

	for _, entry := range TestCases {
		// First encode the MacIface to a bunch of bytes.
		buf, err := EncodeFromMacIface(entry.Mac, entry.Iface)
		assert.Nil(t, err)

		result, err := DecodeToMacIface(buf)
		assert.Nil(t, err)
		assert.Equal(t, entry.Mac, result.Mac)
		assert.Equal(t, entry.Iface, result.Iface)
	}
}

func TestIPFromInterface(t *testing.T) {
	interfaces, err := net.Interfaces()
	assert.Nil(t, err)

	// We can't actually enforce that we get a valid IP, but either the error
	// or the pointer should be nil.
	for _, i := range interfaces {
		addr, err := IPFromInterface(i.Name)
		if err == nil {
			assert.NotNil(t, addr)
		} else {
			assert.Nil(t, addr)
		}
	}
}

func TestIPFromInterfaceNegative(t *testing.T) {
	// Test some fake interfaces.
	var NegativeTestCases = []struct {
		iface string
	}{
		{"fake-interface-0"},
		{"fake-interface-1"},
	}

	for _, tc := range NegativeTestCases {
		addr, err := IPFromInterface(tc.iface)
		assert.Nil(t, addr)
		assert.NotNil(t, err)
	}
}
