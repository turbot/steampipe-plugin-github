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
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Parameters Parameters `json:"parameters"`
}

type Parameters struct {
	Type                               string                             `graphql:"__typename"`
	PullRequestParameters              PullRequestParameters              `graphql:"... on PullRequestParameters"`
	CodeScanningParameters             CodeScanningParameters             `graphql:"... on CodeScanningParameters"`
	CommitAuthorEmailPatternParameters CommitAuthorEmailPatternParameters `graphql:"... on CommitAuthorEmailPatternParameters"`
	CommitMessagePatternParameters     CommitMessagePatternParameters     `graphql:"... on CommitMessagePatternParameters"`
	CommitterEmailPatternParameters    CommitterEmailPatternParameters    `graphql:"... on CommitterEmailPatternParameters"`
	FileExtensionRestrictionParameters FileExtensionRestrictionParameters `graphql:"... on FileExtensionRestrictionParameters"`
	FilePathRestrictionParameters      FilePathRestrictionParameters      `graphql:"... on FilePathRestrictionParameters"`
	MaxFilePathLengthParameters        MaxFilePathLengthParameters        `graphql:"... on MaxFilePathLengthParameters"`
	MaxFileSizeParameters              MaxFileSizeParameters              `graphql:"... on MaxFileSizeParameters"`
	RequiredDeploymentsParameters      RequiredDeploymentsParameters      `graphql:"... on RequiredDeploymentsParameters"`
	RequiredStatusChecksParameters     RequiredStatusChecksParameters     `graphql:"... on RequiredStatusChecksParameters"`
	TagNamePatternParameters           TagNamePatternParameters           `graphql:"... on TagNamePatternParameters"`
	UpdateParameters                   UpdateParameters                   `graphql:"... on UpdateParameters"`
	WorkflowsParameters                WorkflowsParameters                `graphql:"... on WorkflowsParameters"`
}

type PullRequestParameters struct {
	DismissStaleReviewsOnPush      bool `json:"dismiss_stale_reviews_on_push"`
	RequireCodeOwnerReview         bool `json:"require_code_owner_review"`
	RequireLastPushApproval        bool `json:"require_last_push_approval"`
	RequiredApprovingReviewCount   int  `json:"required_approving_review_count"`
	RequiredReviewThreadResolution bool `json:"required_review_thread_resolution"`
}

type CodeScanningParameters struct {
	CodeScanningTools []CodeScanningTool `json:"code_scanning_tools"`
}

type CodeScanningTool struct {
	AlertsThreshold         string `json:"alerts_threshold"`
	SecurityAlertsThreshold string `json:"security_alerts_threshold"`
	Tool                    string `json:"tool"`
}

type CommitAuthorEmailPatternParameters struct {
	Name     string `json:"name"`
	Negate   bool   `json:"negate"`
	Operator string `json:"operator"`
	Pattern  string `json:"pattern"`
}

type CommitMessagePatternParameters struct {
	Name     string `json:"name"`
	Negate   bool   `json:"negate"`
	Operator string `json:"operator"`
	Pattern  string `json:"pattern"`
}

type CommitterEmailPatternParameters struct {
	Name     string `json:"name"`
	Negate   bool   `json:"negate"`
	Operator string `json:"operator"`
	Pattern  string `json:"pattern"`
}

type FileExtensionRestrictionParameters struct {
	RestrictedFileExtensions []string `json:"restricted_file_extensions"`
}

type FilePathRestrictionParameters struct {
	RestrictedFilePaths []string `json:"restricted_file_paths"`
}

type MaxFilePathLengthParameters struct {
	MaxFilePathLength int `json:"max_file_path_length"`
}

type MaxFileSizeParameters struct {
	MaxFileSize int `json:"max_file_size"`
}

type RequiredDeploymentsParameters struct {
	RequiredDeploymentEnvironments []string `json:"required_deployment_environments"`
}

type RequiredStatusChecksParameters struct {
	RequiredStatusChecks             []StatusCheckConfiguration `json:"required_status_checks"`
	StrictRequiredStatusChecksPolicy bool                       `json:"strict_required_status_checks_policy"`
}

type StatusCheckConfiguration struct {
	Context       string `json:"context"`
	IntegrationId int    `json:"integration_id"`
}

type TagNamePatternParameters struct {
	Name     string `json:"name"`
	Negate   bool   `json:"negate"`
	Operator string `json:"operator"`
	Pattern  string `json:"pattern"`
}

type UpdateParameters struct {
	UpdateAllowsFetchAndMerge bool `json:"update_allows_fetch_and_merge"`
}

type WorkflowsParameters struct {
	Workflows []WorkflowFileReference `json:"workflows"`
}

type WorkflowFileReference struct {
	Path         string `json:"path"`
	Ref          string `json:"ref"`
	RepositoryId int    `json:"repository_id"`
	Sha          string `json:"sha"`
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
