package main

import (
	"testing"
)

func TestConnectAndRun(t *testing.T) {
	b := ConnectAndRun("localhost", 22, "chiao", "anna1205", "whoami")
	t.Log(b.String())
}
