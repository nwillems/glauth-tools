package chpass

import (
	"fmt"
	"testing"
)

/*
config.cfg:
[ldaps]
  enabled = false
[ldap]
  enabled = true
  listen = "0.0.0.0:9389"
[backend]
  datastore = "config"
  baseDN = "dc=example,dc=com"
[[users]]
  name = "nwillems"
  mail = "jdoe@example.com"
  unixid = 5001
  primarygroup = 5501
  passsha256 = "6478579e37aff45f013e14eeb30b3cc56c72ccdc310123bcdf53e0333e3f416a" # dogood
[[groups]]
  name = "users"
  unixid = 5501

docker run -v ${PWD}/config.cfg:/app/config/config.cfg -p 9389:9389 glauth/glauth
*/

func TestLdapBind(t *testing.T) {
	username := "nwillems"
	basedn := "ou=users,dc=example,dc=com"

	password := "dogood"
	server := "ldap://127.0.0.1"

	res, err := Bind(username, password, server, basedn)

	if err != nil {
		t.Error(err)
	}

	if res != true {
		fmt.Println("Username/Password is wrong")
		t.Fail()
	}
}
