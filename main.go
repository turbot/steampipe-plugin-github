package main

import (
	"github.com/turbot/steampipe-plugin-github/github"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: github.Plugin})
}
