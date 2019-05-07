package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// https://www.reddit.com/r/golang/comments/87hi86/interactive_ssh/
// Execute a shell and pipe in commands into stdin
// https://studygolang.com/articles/7675

func main() {
	test := make(chan int, 3)
	// go func() {
	// 	i := 0
	// 	select {
	// 	case test <- i:
	// 		i++
	// 	}
	// }()
	test <- 0
	test <- 1
	test <- 3
	// fmt.Println(<-test, <-test)
	go func() {
		for {
			select {
			case i := <-test:
				fmt.Printf("%d\t", i)
			}
		}
	}()
	test <- 4
	fmt.Println()
	os.Exit(0)

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

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		log.Fatal(err)
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	e, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}

	in, out := MuxShell(w, r, e)
	in <- "pwd"
	in <- "pwd"
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}
	// <-out //ignore the shell output

	in <- "exit"
	in <- "exit"

	fmt.Printf("%s\n%s\n", <-out, <-out)

	_, _ = <-out, <-out
	session.Wait()
}
func MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 3)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.Contains(result, "Username:") ||
				strings.Contains(result, "Password:") ||
				strings.Contains(result, "#") {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
