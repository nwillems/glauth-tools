package chpass

import (
	"strings"

	"github.com/BurntSushi/toml"
	glauth "github.com/glauth/glauth/pkg/config"
)

// ParseGLAuthConfig parses a set of configuration according to the same method as used by GLAuth
func ParseGLAuthConfig(data string) (glauth.Config, error) {
	cfg := glauth.Config{}
	// setup defaults
	cfg.LDAP.Enabled = false
	cfg.LDAPS.Enabled = true
	cfg.Backend.NameFormat = "cn"
	cfg.Backend.GroupFormat = "ou"
	cfg.Backend.SSHKeyAttr = "sshPublicKey"

	if _, err := toml.Decode(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

//DeserializeGLAuthConfig takes a configuration struct and encodes it as glauth expects their config file
func DeserializeGLAuthConfig(config glauth.Config) (string, error) {
	writer := strings.Builder{}
	err := (toml.NewEncoder(&writer)).Encode(config)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}
