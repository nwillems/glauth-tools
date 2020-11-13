package cmds

import (
	"testing"

	"github.com/nwillems/glauth-tools/mails"
)

func TestSendEmail(t *testing.T) {
	config := mails.Mailconfig{
		SmtpFromAddress:    "you@example.com",
		SmtpHostName:       "127.0.0.1",
		SmtpHostPort:       9025,
		SmtpPassword:       "herp",
		SmtpUser:           "snerp",
		StartTLS:           false,
		VerifyCertificates: false,
	}

	emailTo := "john@example.com" // Change to something you receive - or mock an smtp server
	emailBody := "Hej med dig, her er en test mail fra golang"
	SendEmail(&emailTo, &emailBody, config)
}
