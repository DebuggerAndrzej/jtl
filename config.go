package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type Config struct {
	Username      string
	Password      string
	Jira_base_url string
	Issues        string
}

func get_toml_config() *Config {
	var config Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	if _, err := toml.DecodeFile(path.Join(homeDir, ".config/jtl.toml"), &config); err != nil {
		fmt.Println(err)
	}

	return &config
}
