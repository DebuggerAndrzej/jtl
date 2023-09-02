package backend

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

func GetTomlConfig(configPath string) *entities.Config {
	var config entities.Config

	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Couldn't determing user's home dir!")
		}
		configPath = path.Join(homeDir, ".config/jtl.toml")
	}

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		panic(fmt.Sprintf("Couldn't load config file. Check if %s file exists!", configPath))
	}

	return &config
}
