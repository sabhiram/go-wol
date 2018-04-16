package main

////////////////////////////////////////////////////////////////////////////////

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////////////////////////////

func TestIPFromInterface(t *testing.T) {
	interfaces, err := net.Interfaces()
	assert.Nil(t, err)

	// We can't actually enforce that we get a valid IP, but either the error
	// or the pointer should be nil.
	for _, i := range interfaces {
		addr, err := ipFromInterface(i.Name)
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
		addr, err := ipFromInterface(tc.iface)
		assert.Nil(t, addr)
		assert.NotNil(t, err)
	}
}
