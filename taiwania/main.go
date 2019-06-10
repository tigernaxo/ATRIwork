package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	user := os.Args[1]
	pw := os.Args[2]
	address := os.Args[3]
	port, err := strconv.Atoi(os.Args[4])
	logErr(err)
	// var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.

	cmd := "cd /home/naxo/文件/gossh; ls -al"
	b, e, err := ConnectAndRun(address, port, user, pw, cmd)
	if b != nil {
		fmt.Print(b.String())
	}
	if e != nil {
		fmt.Print(e.String())
	}
}
