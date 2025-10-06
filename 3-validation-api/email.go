package main

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(to, hash string) error {
	cfg := LoadConfig()

	e := email.NewEmail()
	e.From = cfg.Email
	e.To = []string{to}
	e.Subject = "Email Verification"
	verifyLink := fmt.Sprintf("http://%s:%d/verify/%s", cfg.Address, cfg.Port, hash)
	e.Text = []byte(fmt.Sprintf("Please verify your email by clicking: %s", verifyLink))

	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, "smtp.gmail.com")
	return e.Send("smtp.gmail.com:587", auth)
}
