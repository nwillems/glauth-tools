package cmds

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"

	"github.com/nwillems/glauth-tools/chpass"
	"github.com/nwillems/glauth-tools/mails"
)

//NewUserCreate creates a new password for a user, sends them an email about it, and prints the password-hash
func NewUserCreate(config mails.Mailconfig, args []string) {
	flags := flag.NewFlagSet("new-user", flag.ExitOnError)
	username := flags.String("username", "", "The username to create the password for")
	email := flags.String("email", "", "Email which to send the generated password to")

	// Parse flags
	flags.Parse(args)

	password := generatePassword()
	passSha := chpass.HashPassword(password)
	body := createEmailBody(username, password)

	SendEmail(email, body, config)

	fmt.Printf("name=\"%s\"\nmail=\"%s\"\npasssha256 = \"%s\"\n", *username, *email, *passSha)
}

func generatePassword() *string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	numLetters := big.NewInt(int64(len(letters)))

	passwordLength := 12 // This should probably be configureable
	passRaw := make([]rune, passwordLength)
	for i := range passRaw {
		n, _ := rand.Int(rand.Reader, numLetters)
		n2 := n.Int64()
		passRaw[i] = letters[n2]
	}

	res := string(passRaw)
	return &res
}

func createEmailBody(username, password *string) *string {
	template := "Din bruger bliver oprettet lige om lidt.\nbrugernavn: %s\npassword: %s\n"
	body := fmt.Sprintf(template, *username, *password)
	return &body
}

// SendEmail ensure that the given recipient will recieve an email with the spceified body
func SendEmail(emailTo, body *string, conf mails.Mailconfig) (bool, error) {
	conf.StartTLS = true
	conf.VerifyCertificates = false

	emailFrom := conf.SmtpFromAddress

	subject := "Ny bruger!\n"

	err := mails.SendMail(&emailFrom, emailTo, &subject, body, conf)
	return (err == nil), err
}
