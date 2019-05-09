package main

import (
	"io"

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
