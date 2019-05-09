package main

import (
	"strconv"

	"golang.org/x/crypto/ssh"
)

type myConf struct {
	Term    *termConf
	Client  *ssh.ClientConfig
	Address string
	Port    int
}

func newTestMyConf(user, pw, address string, port int) *myConf {
	return &myConf{
		Term:    defaultTermConf(),
		Client:  testClientConf(user, pw),
		Address: address,
		Port:    port,
	}
}

func (c *myConf) getPortedAddress() string {
	return c.Address + ":" + strconv.Itoa(c.Port)
}

type termConf struct {
	Type string            // term
	H    int               // 24
	W    int               // 80
	Mode ssh.TerminalModes // ssh.TerminalModes{}
}

func defaultTermConf() *termConf {
	return &termConf{
		Type: "term",
		H:    24,
		W:    80,
		Mode: ssh.TerminalModes{},
	}
}

// How To Stop ?
func tcpSSH(c *myConf, inChan <-chan string, outChan, errChan chan<- string) {
	client, err := ssh.Dial("tcp", c.getPortedAddress(), c.Client)
	logErr(err)
	defer client.Close()

	session, err := client.NewSession()
	logErr(err)
	defer session.Close()

	err = session.RequestPty(c.Term.Type, c.Term.H, c.Term.W, c.Term.Mode)
	logErr(err)

	// Get terminal stdin, stdout, stderr
	i, o, e := getPipes(session)
	// Write to session stdin while cmd<-inChan
	// Read session stdout, stderr to outChan, errChan
	go chanToWriter(inChan, i)
	go readerToChan(outChan, o)
	go readerToChan(errChan, e)

	// Start session shell than Wait
	err = session.Shell()
	logErr(err)
	err = session.Wait()
	logErr(err)
}
