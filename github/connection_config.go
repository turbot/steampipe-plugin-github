package github

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type githubConfig struct {
	Token          *string `hcl:"token,optional"`
	BaseURL        *string `hcl:"base_url,optional"`
	AppId          *string `hcl:"app_id,optional"`
	InstallationId *string `hcl:"app_installation_id,optional"`
	PrivateKey     *string `hcl:"app_private_key,optional"`
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
