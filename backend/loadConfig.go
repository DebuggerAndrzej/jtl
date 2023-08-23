package backend

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

func GetTomlConfig() *entities.Config {
	var config entities.Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Couldn't determing user's home dir!")
	}
	if _, err := toml.DecodeFile(path.Join(homeDir, ".config/jtl.toml"), &config); err != nil {
		panic("Couldn't load config file. Check if  ~/.config/jtl.toml file exists!")
	}

	return &config
}
