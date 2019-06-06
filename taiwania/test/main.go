package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

// https://www.reddit.com/r/golang/comments/87hi86/interactive_ssh/
// Execute a shell and pipe in commands into stdin
// https://studygolang.com/articles/7675

func main() {
	// Set pty size
	terminalHeight := 24
	terminalWidth := 80

	// Set login information
	address := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	logErr(err)
	user := os.Args[3]
	var pw string
	if len(os.Args) >= 5 {
		pw = os.Args[4]
	} else {
		pw = ""
	}

	// Set a config
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Dial to server
	client, err := ssh.Dial("tcp", address+":"+strconv.Itoa(port), sshConfig)
	logErr(err)
	defer client.Close()

	// Get client session
	session, err := client.NewSession()
	logErr(err)
	defer session.Close()

	// Request pty using session
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}
	err = session.RequestPty(termType, terminalHeight, terminalWidth, ssh.TerminalModes{})
	logErr(err)

	// Get session pipes
	i, o, e := getPipes(session)
	inChan := make(chan string, 10)
	go chanToWriter(inChan, i)
	outChan, errChan := make(chan []byte, 10), make(chan []byte, 10)
	go readerToChan(outChan, o)
	go readerToChan(errChan, e)

	// 測試從Channel提取server stdout, stderr
	// 這裡可以接log file，讓使用者存取? 如果檔案太大怎麼辦？ 壓縮？
	// 決定先不讓使用者存取，讓使用者存取輸出的檔案(ex: log)就好
	// 直接在登入節點上傳任務資料夾到google drive
	go func() {
		for {
			select {
			case out := <-outChan:
				fmt.Printf("[outChan] %s", string(out))
			case err := <-errChan:
				fmt.Printf("[errChan] %s", string(err))
			}
		}
	}()

	// 測試送指令，這裡可以接到local web api, 讓web送指定
	go func() {
		time.Sleep(time.Second)
		inChan <- "cd /home/naxo/文件/gossh"
		time.Sleep(time.Second)
		inChan <- "ls"
		time.Sleep(time.Second)
		inChan <- "exit"
	}()

	err = session.Shell()
	logErr(err)
	err = session.Wait()
	logErr(err)
}
