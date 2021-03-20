package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAllCommands(t *testing.T) {
	assert.NotEmpty(t, getAllCommands())
}

func Test_getAllOptions(t *testing.T) {
	assert.NotEmpty(t, getAllOptions())
}

func Test_getAppUsageString(t *testing.T) {
	assert.NotEmpty(t, getAppUsageString())
}
