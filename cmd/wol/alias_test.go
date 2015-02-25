package main

import (
	"bytes"
	"encoding/gob"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"os"
	"regexp"
	"runtime"
	"testing"
)

// Helper regex to strip the preamble from the function name. This is
// used to create a temp db file per test (based on the test name).
var RE_stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)

// Validate the DecodeToMacIface function
func TestDecodeToMacIface(test *testing.T) {
	var TestCases = []MacIface{
		{"00:00:00:00:00:00", ""},
		{"00:00:00:00:00:AA", "eth1"},
	}

	for _, entry := range TestCases {
		// First encode the MacIface to a bunch of bytes
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(entry)
		assert.Nil(test, err)

		// Invoke the function and validate that it is equal
		// to our starting MacIface
		result, err := DecodeToMacIface(buf)
		assert.Nil(test, err)
		assert.Equal(test, entry.Mac, result.Mac)
		assert.Equal(test, entry.Iface, result.Iface)
	}
}

// Validate the EncodeFromMacIface function
func TestEncodeFromMacIface(test *testing.T) {
	var TestCases = []MacIface{
		{"00:00:00:00:00:00", "eth0"},
		{"00:00:00:00:00:AA", ""},
	}

	for _, entry := range TestCases {
		// First encode the MacIface to a bunch of bytes
		buf, err := EncodeFromMacIface(entry.Mac, entry.Iface)
		assert.Nil(test, err)

		// Invoke the function and validate that it is equal
		// to our starting MacIface
		result, err := DecodeToMacIface(buf)
		assert.Nil(test, err)
		assert.Equal(test, entry.Mac, result.Mac)
		assert.Equal(test, entry.Iface, result.Iface)
	}
}

// Validate that an invalid db path errors out
func TestInvalidDbPath(test *testing.T) {
	aliases, err := LoadAliases("./dir/no/existy/_test_TestInvalidDbPath")
	assert.NotNil(test, err)
	assert.Nil(test, aliases)
}

//////////////////////////////////////////////////////////////////////////////
// Test suite: AliasDBTests
//      Validate various parts of the DB functionality which needs
//      a common Setup and Teardown to create the db instance
//////////////////////////////////////////////////////////////////////////////
// Define the test suite, and common members which can be accessed from each
// test case within the suite.
type AliasDBTests struct {
	suite.Suite
	dbName  string
	aliases *Aliases
}

// The Setup function is responsible for creating a temporary
// BoltDB to test against. Then, it returns the path to the
// db it creates so we can clean this up at teardown time (along
// with a pointer to the db instance).
func (suite *AliasDBTests) SetupTest() {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		suite.dbName = RE_stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
	}

	var err error
	suite.aliases, err = LoadAliases("./" + suite.dbName)
	assert.Nil(suite.T(), err)
}

// The TearDown function closes the connection to the DB, and
// removes the temporary file created for the same
func (suite *AliasDBTests) TearDownTest() {
	// Close the connection to the bolt db
	err := suite.aliases.Close()
	assert.Nil(suite.T(), err)

	// Remove the temporary test db
	err = os.Remove("./" + suite.dbName)
	assert.Nil(suite.T(), err)
}

// Validates the Aliases Add function
func (suite *AliasDBTests) TestAddAlias() {
	var TestCases = []struct {
		alias, mac, iface string
	}{
		{"one", "00:00:00:00:00:00", "eth0"},
		{"two", "00:00:00:00:00:AA", "eth1"},
		{"thr", "00:00:00:00:11:00", ""},
		{"fou", "00:00:00:00:11:AA", ""},
	}

	entryCount := 0
	for _, entry := range TestCases {
		// Add the alias to the db
		err := suite.aliases.Add(entry.alias, entry.mac, entry.iface)
		assert.Nil(suite.T(), err)
		entryCount += 1

		// Validate that we have "entryCount" number of aliases added
		list, err := suite.aliases.List()
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), entryCount, len(list))

		// Check to ensure that the current map contains the key we
		// just added to the db
		assert.Equal(suite.T(), entry.mac, list[entry.alias].Mac)
		assert.Equal(suite.T(), entry.iface, list[entry.alias].Iface)
	}
}

// Adding a duplicate entry should overwrite the original one
func (suite *AliasDBTests) TestAddDuplicateAlias() {
	err := suite.aliases.Add("test01", "00:11:22:33:44:55", "eth0")
	assert.Nil(suite.T(), err)

	// Validate the first entry exists
	list, err := suite.aliases.List()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "00:11:22:33:44:55", list["test01"].Mac)
	assert.Equal(suite.T(), "eth0", list["test01"].Iface)

	err = suite.aliases.Add("test01", "00:11:22:33:44:66", "")
	assert.Nil(suite.T(), err)

	list, err = suite.aliases.List()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "00:11:22:33:44:66", list["test01"].Mac)
	assert.Equal(suite.T(), "", list["test01"].Iface)
}

// Adding a duplicate entry should overwrite the original one
func (suite *AliasDBTests) TestDeleteAlias() {
	var err error
	var list map[string]MacIface

	err = suite.aliases.Add("test01", "00:11:22:33:44:55", "eth0")
	assert.Nil(suite.T(), err)
	err = suite.aliases.Add("test02", "00:11:22:33:44:66", "")
	assert.Nil(suite.T(), err)

	// Validate that we have two items in the db
	list, err = suite.aliases.List()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(list))

	// Remove test01
	err = suite.aliases.Del("test01")
	assert.Nil(suite.T(), err)
	list, err = suite.aliases.List()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(list))

	// Remove test02
	err = suite.aliases.Del("test02")
	assert.Nil(suite.T(), err)
	list, err = suite.aliases.List()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(list))
}

// Adding a duplicate entry should overwrite the original one
func (suite *AliasDBTests) TestGetAlias() {
	var mi MacIface

	var TestCases = []struct {
		alias, mac, iface string
	}{
		{"one", "00:00:00:00:00:00", "eth0"},
		{"two", "00:00:00:00:00:AA", "eth1"},
		{"thr", "00:00:00:00:11:00", ""},
		{"fou", "00:00:00:00:11:AA", ""},
	}

	for _, entry := range TestCases {
		err := suite.aliases.Add(entry.alias, entry.mac, entry.iface)
		assert.Nil(suite.T(), err)

		mi, err = suite.aliases.Get(entry.alias)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), entry.mac, mi.Mac)
		assert.Equal(suite.T(), entry.iface, mi.Iface)
	}

	// Negative test case - aliases which do not exist
	mi, err := suite.aliases.Get("foobar")
	assert.NotNil(suite.T(), err)
}

// Group up all the test suites we wish to run and dispatch them here
func TestRunAllSuites(t *testing.T) {
	suite.Run(t, new(AliasDBTests))
}
