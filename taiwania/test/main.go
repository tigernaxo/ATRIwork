package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// https://www.reddit.com/r/golang/comments/87hi86/interactive_ssh/
// Execute a shell and pipe in commands into stdin
// https://studygolang.com/articles/7675

func main() {
	terminalHeight := 24
	terminalWidth := 80
	user := os.Args[2]
	address := os.Args[1]
	var pw string
	if len(os.Args) >= 4 {
		pw = os.Args[3]
	} else {
		pw = ""
	}
	if len(os.Args) < 4 {
		fmt.Println("Usage...")
		os.Exit(0)
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", address, sshConfig)
	logErr(err)
	defer client.Close()

	session, err := client.NewSession()
	logErr(err)
	defer session.Close()

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = session.RequestPty(termType, terminalHeight, terminalWidth, ssh.TerminalModes{})
	logErr(err)

	i, o, e := getPipes(session)
	// 可以設置chan取得out跟err給http server使用
	go io.Copy(os.Stderr, e)
	go io.Copy(os.Stdout, o)
	inChan := make(chan string, 10)
	go chanToWriter(inChan, i)

	go func() {
		inChan <- "ll"
		inChan <- "exit"
	}()

	err = session.Shell()
	logErr(err)
	err = session.Wait()
	logErr(err)
}
