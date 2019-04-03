package main

import (
	"testing"
)

func TestConnectAndRun(t *testing.T) {
	b := ConnectAndRunOnce("ptt.cc", 22, "bbsu", "", "")
	t.Log(b.String())
}
