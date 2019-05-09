package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// https://www.reddit.com/r/golang/comments/87hi86/interactive_ssh/
// Execute a shell and pipe in commands into stdin
// https://studygolang.com/articles/7675

func main() {
	terminalHeight := 24
	terminalWidth := 80
	user := os.Args[3]
	address := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	logErr(err)
	var pw string
	if len(os.Args) >= 5 {
		pw = os.Args[4]
	} else {
		pw = ""
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", address+":"+strconv.Itoa(port), sshConfig)
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
	// go io.Copy(os.Stdout, o)
	// outChan := make(chan string, 10)
	// go readerToChan(outChan, o)
	inChan := make(chan string, 10)
	go chanToWriter(inChan, i)

	go func() {
		inChan <- "ll"
		inChan <- "exit"
	}()
	buf, all := make([]byte, 256), make([]byte, 256)

	go func() {
		for {
			n, err := o.Read(buf)
			logErr(err)
			// n永遠不會為0，需要修改readerToChan()
			switch buf[n-1] {
			case 10:
				all = append(all, buf[:n]...)
				fmt.Printf("%v", string(all))
				all = all[:0]
			default:
				all = append(all, buf[:n]...)
			}
		}
	}()

	fmt.Println(o, e)
	err = session.Shell()
	logErr(err)
	err = session.Wait()
	logErr(err)
}
