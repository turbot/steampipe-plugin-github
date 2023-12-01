package models

import "github.com/shurcooL/githubv4"

type Environment struct {
	Id     int    `graphql:"id: databaseId @include(if:$includeEnvironmentId)" json:"id,omitempty"`
	NodeId string `graphql:"nodeId: id @include(if:$includeEnvironmentNodeId)" json:"node_id,omitempty"`
	Name   string `graphql:"name @include(if:$includeEnvironmentName)" json:"name,omitempty"`
	// protectionRules [pageable]
}

type Deployment struct {
	Id                  int                      `graphql:"id: databaseId @include(if:$includeDeploymentId)" json:"id,omitempty"`
	NodeId              string                   `graphql:"nodeId: id @include(if:$includeDeploymentNodeId)" json:"node_id,omitempty"`
	CommitSha           string                   `graphql:"sha: commitOid @include(if:$includeDeploymentCommitSha)" json:"sha"`
	CreatedAt           NullableTime             `graphql:"createdAt @include(if:$includeDeploymentCreatedAt)" json:"created_at,omitempty"`
	Creator             Actor                    `graphql:"creator @include(if:$includeDeploymentCreator)" json:"creator,omitempty"`
	Description         string                   `graphql:"description @include(if:$includeDeploymentDescription)" json:"description,omitempty"`
	Environment         string                   `graphql:"environment @include(if:$includeDeploymentEnvironment)" json:"environment,omitempty"`
	LatestEnvironment   string                   `graphql:"latestEnvironment @include(if:$includeDeploymentLatestEnvironment)" json:"latest_environment,omitempty"`
	LatestStatus        DeploymentStatus         `graphql:"latestStatus @include(if:$includeDeploymentLatestStatus)" json:"latest_status,omitempty"`
	OriginalEnvironment string                   `graphql:"originalEnvironment @include(if:$includeDeploymentOriginalEnvironment)" json:"original_environment,omitempty"`
	Payload             string                   `graphql:"payload @include(if:$includeDeploymentPayload)" json:"payload,omitempty"`
	Ref                 BasicRef                 `graphql:"ref @include(if:$includeDeploymentRef)" json:"ref,omitempty"`
	State               githubv4.DeploymentState `graphql:"state @include(if:$includeDeploymentState)" json:"state,omitempty"`
	Task                string                   `graphql:"task @include(if:$includeDeploymentTask)" json:"task,omitempty"`
	UpdatedAt           NullableTime             `graphql:"updatedAt @include(if:$includeDeploymentUpdatedAt)" json:"updated_at,omitempty"`
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
