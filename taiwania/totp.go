package main

import (
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)

// GetTOTP direct get totp from secret
func GetTOTP(secret string) string {
	otp, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		log.Panic(err)
	}
	return otp
}
