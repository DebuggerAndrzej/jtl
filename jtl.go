package main

import (
	"flag"

	"github.com/DebuggerAndrzej/jtl/backend"
	"github.com/DebuggerAndrzej/jtl/ui"
)

func main() {
	var configPath = flag.String("config", "", "Full path to config file")
	flag.Parse()

	config := backend.GetTomlConfig(*configPath)
	client := backend.GetJiraClient(config)
	ui.InitTui(config, client, *configPath)
}
