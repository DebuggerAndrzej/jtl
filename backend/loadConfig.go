package backend

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"

	"jtl/backend/entities"
)

func GetTomlConfig() *entities.Config {
	var config entities.Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	if _, err := toml.DecodeFile(path.Join(homeDir, ".config/jtl.toml"), &config); err != nil {
		fmt.Println(err)
	}

	return &config
}
