package github

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type githubConfig struct {
	Token          *string `hcl:"token"`
	BaseURL        *string `hcl:"base_url"`
	AppId          *string `hcl:"app_id"`
	InstallationId *string `hcl:"app_installation_id"`
	PrivateKey     *string `hcl:"app_private_key"`
}

func ConfigInstance() interface{} {
	return &githubConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) githubConfig {
	if connection == nil || connection.Config == nil {
		return githubConfig{}
	}
	config, _ := connection.Config.(githubConfig)
	return config
}
