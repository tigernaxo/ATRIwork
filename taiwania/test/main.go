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
	sshConfig := getInsecurePwSSHConfig(user, pw)

	// Dial to server
	client, err := ssh.Dial("tcp", address+":"+strconv.Itoa(port), sshConfig)
	logErr(err)
	defer client.Close()

	// Get client session
	session, err := client.NewSession()
	logErr(err)
	defer session.Close()

	requestDefaultPty(session)

	// Get session pipes
	inChan, outChan, errChan := make(chan string, 10), make(chan []byte, 10), make(chan []byte, 10)
	pipeToChan(session, inChan, outChan, errChan)

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
