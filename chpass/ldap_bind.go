package chpass

import (
	"fmt"

	ldap "github.com/go-ldap/ldap"
)

// Bind binds to the given server with the given credentials, returning the base user info
func Bind(username, password, server, basedn string) (bool, error) {
	conn, err := ldap.DialURL(server)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	usernameDN := fmt.Sprintf("cn=%s,%s", username, basedn)
	err = conn.Bind(usernameDN, password)
	if err != nil {
		return false, err
	}
	return true, nil
}
