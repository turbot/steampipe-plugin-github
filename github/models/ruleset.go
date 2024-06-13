package models

type Ruleset struct {
	CreatedAt    string        `json:"created_at"`
	DatabaseID   int           `json:"database_id"`
	Enforcement  string        `json:"enforcement"`
	Name         string        `json:"name"`
	ID           string        `json:"id"`
	Rules        []Rule        `json:"rules"`
	BypassActors []BypassActor `json:"bypass_actors"`
	Conditions   Conditions    `json:"conditions"`
}

type Rule struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

type BypassActor struct {
	BypassMode               string `json:"bypass_mode"`
	DeployKey                bool   `json:"deploy_key"`
	ID                       string `json:"id"`
	RepositoryRoleDatabaseID int    `json:"repository_role_database_id"`
	RepositoryRoleName       string `json:"repository_role_name"`
}

type Conditions struct {
	RefName struct {
		Exclude []string `json:"exclude"`
		Include []string `json:"include"`
	} `json:"ref_name"`
	RepositoryID struct {
		RepositoryIds []string `json:"repository_ids"`
	} `json:"repository_id"`
	RepositoryName struct {
		Exclude   []string `json:"exclude"`
		Include   []string `json:"include"`
		Protected bool     `json:"protected"`
	} `json:"repository_name"`
}
