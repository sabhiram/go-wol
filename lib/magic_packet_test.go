package wol

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetMacBytes(test *testing.T) {
    var PositiveTestCases = []struct {
        mac      string
        expected MacAddress
    } {
        { "00:00:00:00:00:00",   MacAddress{0,0,0,0,0,0} },
        { "00:00:00:00:00:00",   MacAddress{0,0,0,0,0,0} },
    }

    for _, t := range PositiveTestCases {
        macAddress, err := GetMacBytes(t.mac)
        assert.Equal(test, t.expected, *macAddress)
        assert.Equal(test, err, nil)
    }
}

func TestGetMacBytesNegative(test *testing.T) {
    var NegativeTestCases = []struct {
        mac   string
    } {
        { "00x00:00:00:00:00" },
        { "00:00:Z0:00:00:00" },
    }

    for _, t := range NegativeTestCases {
        _, err := GetMacBytes(t.mac)
        assert.NotEqual(test, err, nil)
    }
}

func TestNewMagicPacket(test *testing.T) {
    var PositiveTestCases = []struct {
        mac      string
        expected MacAddress
    } {
        { "00:00:00:00:00:00",   MacAddress{0,0,0,0,0,0} },
        { "00:00:00:00:00:00",   MacAddress{0,0,0,0,0,0} },
    }

    for _, t := range PositiveTestCases {
        magicPkt, err := NewMagicPacket(t.mac)
        assert.Equal(test, err, nil)
    }
}
