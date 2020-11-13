package mails

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Mailconfig struct {
	SmtpHostName       string
	SmtpHostPort       string
	SmtpUser           string
	SmtpPassword       string
	SmtpFromAddress    string
	StartTLS           bool
	VerifyCertificates bool
}

// SendMail, mirrors the sendmail function from the smtp package. However it takes only a single
//  recipient, a mail-configuration and doesn't do verifications.
func SendMail(from, to, subject, body *string, conf Mailconfig) error {
	addr := fmt.Sprintf("%s:%s", conf.SmtpHostName, conf.SmtpHostPort)
	dataBody := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", *from, *to, *subject, *body)
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Hello("glauth-tools.local"); err != nil {
		return err
	}

	if hasTLS, _ := client.Extension("STARTTLS"); conf.StartTLS && hasTLS {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: !conf.VerifyCertificates,
			ServerName:         conf.SmtpHostName,
		}

		client.StartTLS(tlsconfig)
	}

	smtpAuth := smtp.PlainAuth("", conf.SmtpUser, conf.SmtpPassword, conf.SmtpHostName)
	if err = client.Auth(smtpAuth); err != nil {
		return err
	}

	// To && From
	if err = client.Mail(*from); err != nil {
		return err
	}

	if err = client.Rcpt(*to); err != nil {
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(dataBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	if err = client.Quit(); err != nil {
		return err
	}

	return nil
}
