package backend

import (
	"fmt"

	"github.com/amazurki/JTL/backend/entities"
)

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
