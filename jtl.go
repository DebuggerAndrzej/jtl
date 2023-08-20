package main

import (
	"github.com/DebuggerAndrzej/jtl/backend"
	"github.com/DebuggerAndrzej/jtl/ui"
)

func main() {
	config := backend.GetTomlConfig()
	client := backend.GetJiraClient(config)
	ui.InitTui(config, client)

}
