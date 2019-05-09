package main

import (
	"io"
	"log"
)

func logErr(err error) {
	if err != nil {
		log.Panicf("[Error] %v\n", err)
	}
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
