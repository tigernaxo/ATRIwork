package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/ssh"
)

// Todo: function open a new shell pipe
// Note: functions maybe use to run more than one command
// func (s *Session) StderrPipe() (io.Reader, error)
// func (s *Session) StdinPipe() (io.WriteCloser, error)
// func (s *Session) StdoutPipe() (io.Reader, error)
// func (s *Session) Wait() error
// http://blog.ralch.com/tutorial/golang-ssh-connection/

// ConnectAndRunOnce direct connect to server then run command
// modified from official example "Dial"
func ConnectAndRunOnce(host string, port int, id, password, command string) *bytes.Buffer {

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

// GetTOTP direct get totp from secret
func GetTOTP(secret string) string {
	otp, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		log.Panic(err)
	}
	return otp
}
