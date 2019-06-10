package main

import (
	"golang.org/x/crypto/ssh"
)

type myConf struct {
	Address  string
	Port     int
	User     string
	Password string
	TermMode string
	HostKey  string
}

// func newTestMyConf(user, pw, address string, port int) *myConf {
// 	return &myConf{
// 		User:     user,
// 		Password: pw,
// 		Address:  address,
// 		Port:     port,
// 		HostKey:  "",
// 		TermMode: "xterm",
// 	}
// }

// func (c *myConf) getPortedAddress() string {
// 	return c.Address + ":" + strconv.Itoa(c.Port)
// }

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

func getInsecurePwSSHConfig(user, pw string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}
