package main

import (
	"github.com/amazurki/JTL/backend"
	"github.com/amazurki/JTL/ui"
)

func main() {
	config := get_toml_config()
	client := get_jira_client(config)
	initTui(config, client)
}
