package mails

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	smtpHost := "127.0.0.1"
	smtpUser := "user"
	smtpPassword := "password"

	emailFrom := "you@example.com"
	emailTo := "new_user@example.com"
	emailSubject := "Subject of email"
	emailBody := "This is a test"
	conf := Mailconfig{
		SmtpHostName:       smtpHost,
		SmtpHostPort:       "587",
		SmtpUser:           smtpUser,
		SmtpPassword:       smtpPassword,
		SmtpFromAddress:    emailFrom,
		StartTLS:           true,
		VerifyCertificates: false,
	}
	SendMail(&emailFrom, &emailTo, &emailSubject, &emailBody, conf)
}
