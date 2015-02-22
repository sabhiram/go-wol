package main

import (
    "bytes"
    "encoding/gob"

    "github.com/stretchr/testify/assert"
    "testing"
)

// This is just a dummy test to enable cross package coverage
func TestDecodeToMacIface(test *testing.T) {
    var TestCases = []MacIface {
        MacIface{ "00:00:00:00:00:00", "eth0" },
        MacIface{ "00:00:00:00:00:AA", "eth1" },
    }

    for _, entry := range TestCases {
        // First encode the MacIface to a bunch of bytes
        buf := bytes.NewBuffer(nil)
        err := gob.NewEncoder(buf).Encode(entry);
        assert.Equal(test, err, nil)

        // Invoke the function and validate that it is equal
        // to our starting MacIface
        result, err := DecodeToMacIface(buf)
        assert.Equal(test, err, nil)
        assert.Equal(test, entry, result)
    }
}

// TODO: Add BoltDB related alias tests here...
