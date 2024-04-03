package models

type CommunityProfile struct {
	LicenseInfo          BaseLicense             `graphql:"licenseInfo @include(if:$includeCPLicense)" json:"license_info"`
	CodeOfConduct        RepositoryCodeOfConduct `graphql:"codeOfConduct @include(if:$includeCPCodeOfConduct)" json:"code_of_conduct"`
	IssueTemplates       []IssueTemplate         `graphql:"issueTemplates @include(if:$includeCPIssueTemplates)" json:"issue_templates"`
	PullRequestTemplates []PullRequestTemplate   `graphql:"pullRequestTemplates @include(if:$includeCPPullRequestTemplates)" json:"pull_request_templates"`
	// readme
	ReadMeLower struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"readMeLower: object(expression: \"HEAD:readme.md\") @include(if:$includeCPReadme)"`
	ReadMeUpper struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"readMeUpper: object(expression: \"HEAD:README.md\") @include(if:$includeCPReadme)"`
	// contributing
	ContributingLower struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"contributingLower: object(expression: \"HEAD:contributing.md\") @include(if:$includeCPContributing)"`
	ContributingTitle struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"contributingTitle: object(expression: \"HEAD:Contributing.md\") @include(if:$includeCPContributing)"`
	ContributingUpper struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"contributingUpper: object(expression: \"HEAD:CONTRIBUTING.md\") @include(if:$includeCPContributing)"`
	// security
	SecurityLower struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"securityLower: object(expression: \"HEAD:security.md\") @include(if:$includeCPSecurity)"`
	SecurityTitle struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"securityTitle: object(expression: \"HEAD:Security.md\") @include(if:$includeCPSecurity)"`
	SecurityUpper struct {
		Blob Blob `graphql:"... on Blob"`
	} `graphql:"securityUpper: object(expression: \"HEAD:SECURITY.md\") @include(if:$includeCPSecurity)"`
}
