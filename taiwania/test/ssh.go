package main

import (
	"io"

	"golang.org/x/crypto/ssh"
)

func getPipes(s *ssh.Session) (io.WriteCloser, io.Reader, io.Reader) {
	i, err := s.StdinPipe()
	logErr(err)
	o, err := s.StdoutPipe()
	logErr(err)
	e, err := s.StderrPipe()
	logErr(err)
	return i, o, e
}
func readerToChan(c chan<- string, r io.Reader) {
	buf := make([]byte, 128)
	s := ""
	for {
		n, err := r.Read(buf)
		logErr(err)

		switch n {
		case 0:
			c <- s
			s = ""
		default:
			s += string(buf)
		}

	}
}
func chanToWriter(c <-chan string, w io.WriteCloser) {
	for {
		select {
		case cmd := <-c:
			w.Write([]byte(cmd + "\n"))
		}
	}
}
