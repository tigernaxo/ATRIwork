package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// only for test: HostKeyCallback: ssh.InsecureIgnoreHostKey()
func testClientConf(user, pw string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func getPipes(s *ssh.Session) (io.WriteCloser, io.Reader, io.Reader) {
	i, err := s.StdinPipe()
	logErr(err)
	o, err := s.StdoutPipe()
	logErr(err)
	e, err := s.StderrPipe()
	logErr(err)
	return i, o, e
}

func pipeToChan(s *ssh.Session, in chan string, out, err chan []byte) {
	i, o, e := getPipes(s)
	go chanToWriter(in, i)
	go readerToChan(out, o)
	go readerToChan(err, e)
}
func requestDefaultPty(s *ssh.Session) {
	// Set pty size
	terminalHeight := 24
	terminalWidth := 80

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}
	err := s.RequestPty(termType, terminalHeight, terminalWidth, ssh.TerminalModes{})
	logErr(err)
}

// Use https://godoc.org/github.com/buildkite/terminal
// Convert ansi to html code

// ConnectAndRun direct connect to server then run command
// modified from official example "Dial"
func ConnectAndRun(host string, port int, id, password, command string) *bytes.Buffer {

	address := fmt.Sprintf("%s:%d", host, port)
	// var hostKey ssh.PublicKey

	// An SSH client is represented with a ClientConn.
	config := &ssh.ClientConfig{
		User: id,
		// To authenticate with the remote server you must:
		//   1.pass at least one implementation of AuthMethod
		//   2.provide a HostKeyCallback.
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback need to migrate to FixedHostKey for production
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can:
	//   1.Execute a single command by session.Run.
	//   2.Get command returned info by bytes.Buffer
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	return &b
}
