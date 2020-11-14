package chpass

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Config struct {
	LdapServer string
	BaseDN     string

	KubernetesConfigmap string
	KubernetesNamespace string
}

//HashPassword hashes the given password using sha256
func HashPassword(password *string) *string {
	hash := sha256.New()
	hash.Write([]byte(*password))

	res := hex.EncodeToString(hash.Sum(nil))
	return &res
}

// ChangePassword changes the given users password, from oldPassword to newPassword. First verifies oldPassword against given LDAP server
func ChangePassword(username, oldPassword, newPassword string, config Config) error {
	// Perform Bind
	validUser, err := Bind(username, oldPassword, config.LdapServer, config.BaseDN)
	if err != nil {
		return err
	}

	if !validUser {
		return fmt.Errorf("User could not be authenticated")
	}
	// Get configmap
	data, err := GetConfigMap(config.KubernetesConfigmap, config.KubernetesNamespace)
	if err != nil {
		return err
	}
	// Parse glauth config
	glauthConfig, err := ParseGLAuthConfig(data["config.cfg"])
	// Exchange password for new hashed
	hashedNewPassword := HashPassword(&newPassword)
	for idx, user := range glauthConfig.Users {
		if user.Name == username {
			glauthConfig.Users[idx].PassSHA256 = *hashedNewPassword
		}
	}
	// Deserialize config
	updatedConfigData, err := DeserializeGLAuthConfig(glauthConfig)
	if err != nil {
		return err
	}

	data["config.cfg"] = updatedConfigData
	// Save configmap
	err = UpdateConfigMap(config.KubernetesConfigmap, config.KubernetesNamespace, data)
	if err != nil {
		return err
	}

	return nil
}
