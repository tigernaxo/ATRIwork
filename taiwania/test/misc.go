package main

import (
	"io"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)

func logErr(err error) {
	if err != nil {
		log.Panicf("[Error] %v\n", err)
	}
}

func readerToChan(c chan<- []byte, r io.Reader) {
	all, buf := make([]byte, 256), make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if err != nil && err.Error() != "EOF" {
			log.Panicf("[Error] %v\n", err)
		}
		if n == 0 {
			continue
		}

		switch buf[n-1] {
		case 10:
			all = append(all, buf[:n]...)
			c <- all
			all = all[:0]
		default:
			all = append(all, buf[:n]...)
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

// GetTOTP direct get totp from secret
func GetTOTP(secret string) string {
	otp, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		log.Panic(err)
	}
	return otp
}
