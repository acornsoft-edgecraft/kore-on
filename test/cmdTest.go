package main

import (
	"bytes"
	"testing"

	"cube/cmd"

	"github.com/hhkbp2/testify/assert"
)

func Test_ExecuteAnyCommand(t *testing.T) {

	actual := new(bytes.Buffer)
	cmd.Execute()

	expected := "This-is-command-a1"

	assert.Equal(t, actual.String(), expected, "actual is not expected")
}
