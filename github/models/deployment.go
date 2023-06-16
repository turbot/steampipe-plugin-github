package models

import "github.com/shurcooL/githubv4"

type Environment struct {
	basicIdentifiers
	// protectionRules [pageable]
}

type Deployment struct {
	Id                  int                      `graphql:"id: databaseId" json:"id,omitempty"`
	NodeId              string                   `graphql:"nodeId: id" json:"node_id,omitempty"`
	CommitSha           string                   `graphql:"sha: commitOid" json:"sha"`
	CreatedAt           NullableTime             `json:"created_at"`
	Creator             Actor                    `json:"creator"`
	Description         string                   `json:"description"`
	Environment         string                   `json:"environment"`
	LatestEnvironment   string                   `json:"latest_environment"`
	LatestStatus        DeploymentStatus         `json:"latest_status"`
	OriginalEnvironment string                   `json:"original_environment"`
	Payload             string                   `json:"payload"`
	Ref                 BasicRef                 `json:"ref"`
	State               githubv4.DeploymentState `json:"state"`
	Task                string                   `json:"task"`
	UpdatedAt           NullableTime             `json:"updated_at"`
}

type DeploymentStatus struct {
	NodeId         string                         `graphql:"nodeId: id" json:"node_id,omitempty"`
	CreatedAt      NullableTime                   `json:"created_at"`
	Creator        Actor                          `json:"creator"`
	Description    string                         `json:"description"`
	EnvironmentUrl string                         `json:"environment_url"`
	LogUrl         string                         `json:"log_url"`
	State          githubv4.DeploymentStatusState `json:"state"`
	UpdatedAt      NullableTime                   `json:"updated_at"`
}
