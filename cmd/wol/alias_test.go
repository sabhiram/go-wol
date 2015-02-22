package main

import (
	"bytes"
	"encoding/gob"

	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// This is just a dummy test to enable cross package coverage
func TestDecodeToMacIface(test *testing.T) {
	var TestCases = []MacIface{
		{"00:00:00:00:00:00", "eth0"},
		{"00:00:00:00:00:AA", "eth1"},
	}

	for _, entry := range TestCases {
		// First encode the MacIface to a bunch of bytes
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(entry)
		assert.Equal(test, nil, err)

		// Invoke the function and validate that it is equal
		// to our starting MacIface
		result, err := DecodeToMacIface(buf)
		assert.Equal(test, nil, err)
		assert.Equal(test, entry.Mac, result.Mac)
		assert.Equal(test, entry.Iface, result.Iface)
	}
}

// The Setup function is responsible for creating a temporary
// BoltDB to test against. Then, it returns the path to the
// db it creates so we can clean this up at teardown time (along
// with a pointer to the db instance).
func AliasSetup(test *testing.T, dbPath string) (*Aliases, string) {
	aliases, err := LoadAliases(dbPath)
	assert.Equal(test, nil, err)
	return aliases, dbPath
}

// The Teardown function closes the connection to the DB, and
// removes the temporary file created for the same
func AliasTeardown(test *testing.T, aliases *Aliases, dbPath string) {
	// Close the connection to the bolt db
	err := aliases.Close()
	assert.Equal(test, nil, err)
	// Remove the temporary test db
	err = os.Remove(dbPath)
	assert.Equal(test, nil, err)
}

func TestAddAlias(test *testing.T) {
	var TestCases = []struct {
		alias, mac, iface string
	}{
		{"one", "00:00:00:00:00:00", "eth0"},
		{"two", "00:00:00:00:00:AA", "eth1"},
		{"three", "00:00:00:00:11:00", ""},
		{"four", "00:00:00:00:11:AA", ""},
	}

	// Open a test db before we start any of the tests
	aliases, dbPath := AliasSetup(test, "./_test_TestAddAlias")
	// Cleanup after the tests run
	defer AliasTeardown(test, aliases, dbPath)

	entryCount := 0
	for _, entry := range TestCases {
		// Add the alias to the db
		err := aliases.Add(entry.alias, entry.mac, entry.iface)
		assert.Equal(test, nil, err)
		entryCount += 1

		// Validate that we have "entryCount" number of aliases added
		list, err := aliases.List()
		assert.Equal(test, nil, err)
		assert.Equal(test, entryCount, len(list))

		// Check to ensure that the current map contains the key we
		// just added to the db
		assert.Equal(test, entry.mac, list[entry.alias].Mac)
		assert.Equal(test, entry.iface, list[entry.alias].Iface)
	}
}

func TestAddDuplicateAlias(test *testing.T) {
	// Open a test db before we start any of the tests
	aliases, dbPath := AliasSetup(test, "./_test_TestAddDuplicateAlias")
	// Cleanup after the tests run
	defer AliasTeardown(test, aliases, dbPath)

	err := aliases.Add("test01", "00:11:22:33:44:55", "eth0")
	assert.Equal(test, nil, err)

	// Validate the first entry exists
	list, err := aliases.List()
	assert.Equal(test, nil, err)
	assert.Equal(test, "00:11:22:33:44:55", list["test01"].Mac)
	assert.Equal(test, "eth0", list["test01"].Iface)

	err = aliases.Add("test01", "00:11:22:33:44:66", "")
	assert.Equal(test, nil, err)

	list, err = aliases.List()
	assert.Equal(test, nil, err)
	assert.Equal(test, "00:11:22:33:44:66", list["test01"].Mac)
	assert.Equal(test, "", list["test01"].Iface)
}
