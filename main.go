package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nwillems/glauth-tools/chpass"
	"github.com/nwillems/glauth-tools/cmds"
	"github.com/nwillems/glauth-tools/mails"
)

const help = `Usage: glauth-tools [command] [--options]
Supported commands:
  new-user - creates a random passwords and emails it to the new user, outputs the password sha for glauth config
  server - starts a password change server
`

type Configuration struct {
	NewUser        mails.Mailconfig
	PasswordChange chpass.Config
}

func readConfiguration(configFilepath string) (Configuration, error) {
	cfg := Configuration{} // Consider if there should be default value

	if _, err := toml.DecodeFile(configFilepath, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	configFilepath, isSetConfigFilepath := os.LookupEnv("GLAUTH_TOOLS_CONFIG")
	if !isSetConfigFilepath {
		configFilepath = "./config.toml"
	}

	configuration, err := readConfiguration(configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[1]
	fmt.Printf("=== %s ===\n", cmd)
	args := os.Args[2:]

	if cmd == "new-user" {
		cmds.NewUserCreate(configuration.NewUser, args)
	} else if cmd == "server" {
		cmds.PasswordChangeServer(configuration.PasswordChange, args)
	} else {
		fmt.Print(help)
	}
}
