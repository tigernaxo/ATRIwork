package main

import (
        "fmt"
        "github.com/pquerna/otp/totp"
        "os"
        "time"
)

func main() {
        secret := os.Args[1]
        code, _ := totp.GenerateCode(secret, time.Now().UTC())
        fmt.Println(code)
