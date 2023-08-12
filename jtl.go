package main

import (
	"jtl/backend"
	"jtl/ui"
)

func main() {
	config := backend.GetTomlConfig()
	client := backend.GetJiraClient(config)
	ui.InitTui(config, client)

}
