package main

import (
	"bytes"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/ssh"
)

// ConnectAndRun direct connect to server then run command
// modified from official example "Dial"
func ConnectAndRun(ip string, port int, id, password, command string) *bytes.Buffer {
	// var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: id,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	client, err := ssh.Dial("tcp", ip+":"+string(port), config)
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

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
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
