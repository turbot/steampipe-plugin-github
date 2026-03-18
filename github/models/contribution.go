package models

import "github.com/shurcooL/githubv4"

type ContributionCalendar struct {
	TotalContributions int                 `graphql:"totalContributions" json:"total_contributions"`
	Weeks              []ContributionWeek  `graphql:"weeks" json:"weeks"`
}

type ContributionWeek struct {
	ContributionDays []ContributionDay `graphql:"contributionDays" json:"contribution_days"`
	FirstDay         githubv4.Date     `graphql:"firstDay" json:"first_day"`
}

type ContributionDay struct {
	Color             string                      `graphql:"color" json:"color"`
	ContributionCount int                         `graphql:"contributionCount" json:"contribution_count"`
	ContributionLevel githubv4.ContributionLevel  `graphql:"contributionLevel" json:"contribution_level"`
	Date              githubv4.Date               `graphql:"date" json:"date"`
	Weekday           int                         `graphql:"weekday" json:"weekday"`
}

type CommitContributionsByRepository struct {
	Repository struct {
		NameWithOwner string `graphql:"nameWithOwner" json:"name_with_owner"`
		Url           string `graphql:"url" json:"url"`
	} `graphql:"repository" json:"repository"`
	Contributions struct {
		TotalCount int `graphql:"totalCount" json:"total_count"`
	} `graphql:"contributions" json:"contributions"`
}

type ContributionsCollection struct {
	TotalCommitContributions                  int `graphql:"totalCommitContributions" json:"total_commit_contributions"`
	TotalIssueContributions                   int `graphql:"totalIssueContributions" json:"total_issue_contributions"`
	TotalPullRequestContributions             int `graphql:"totalPullRequestContributions" json:"total_pull_request_contributions"`
	TotalPullRequestReviewContributions       int `graphql:"totalPullRequestReviewContributions" json:"total_pull_request_review_contributions"`
	TotalRepositoriesWithContributedCommits   int `graphql:"totalRepositoriesWithContributedCommits" json:"total_repositories_with_contributed_commits"`
	ContributionCalendar                      ContributionCalendar             `graphql:"contributionCalendar @include(if:$includeContributionCalendar)" json:"contribution_calendar,omitempty"`
	CommitContributionsByRepository           []CommitContributionsByRepository `graphql:"commitContributionsByRepository(maxRepositories: $maxRepositories) @include(if:$includeCommitContributionsByRepository)" json:"commit_contributions_by_repository,omitempty"`
}
