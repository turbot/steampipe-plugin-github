package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubUserContributionStats() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user_contribution_stats",
		Description: "Contribution summary and calendar data for a GitHub user.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "login", Require: plugin.Required},
				{Name: "from_date", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "to_date", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "max_repositories", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubUserContributionStatsList,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user.", Transform: transform.FromQual("login")},
			{Name: "from_date", Type: proto.ColumnType_TIMESTAMP, Description: "Start date for the contribution window.", Transform: transform.FromQual("from_date")},
			{Name: "to_date", Type: proto.ColumnType_TIMESTAMP, Description: "End date for the contribution window.", Transform: transform.FromQual("to_date")},
			{Name: "max_repositories", Type: proto.ColumnType_INT, Description: "Maximum repositories returned for commit contributions by repository.", Transform: transform.FromQual("max_repositories")},
			{Name: "total_commit_contributions", Type: proto.ColumnType_INT, Description: "Total count of commit contributions.", Transform: transform.FromField("TotalCommitContributions")},
			{Name: "total_issue_contributions", Type: proto.ColumnType_INT, Description: "Total count of issue contributions.", Transform: transform.FromField("TotalIssueContributions")},
			{Name: "total_pull_request_contributions", Type: proto.ColumnType_INT, Description: "Total count of pull request contributions.", Transform: transform.FromField("TotalPullRequestContributions")},
			{Name: "total_pull_request_review_contributions", Type: proto.ColumnType_INT, Description: "Total count of pull request review contributions.", Transform: transform.FromField("TotalPullRequestReviewContributions")},
			{Name: "total_repositories_with_contributed_commits", Type: proto.ColumnType_INT, Description: "Total count of repositories with contributed commits.", Transform: transform.FromField("TotalRepositoriesWithContributedCommits")},
			{Name: "contribution_calendar", Type: proto.ColumnType_JSON, Description: "Contribution calendar with weeks and days.", Hydrate: contributionHydrateCalendar, Transform: transform.FromValue().NullIfZero()},
			{Name: "commit_contributions_by_repository", Type: proto.ColumnType_JSON, Description: "Commit contributions aggregated by repository.", Hydrate: contributionHydrateCommitContributionsByRepository, Transform: transform.FromValue().NullIfZero()},
		}),
	}
}

func tableGitHubUserContributionStatsList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	login := d.EqualsQuals["login"].GetStringValue()

	var fromDate *githubv4.DateTime
	var toDate *githubv4.DateTime
	maxRepositories := 100

	if d.EqualsQuals["from_date"] != nil {
		fromTime := d.EqualsQuals["from_date"].GetTimestampValue().AsTime()
		fromDate = githubv4.NewDateTime(githubv4.DateTime{Time: fromTime})
	}

	if d.EqualsQuals["to_date"] != nil {
		toTime := d.EqualsQuals["to_date"].GetTimestampValue().AsTime()
		toDate = githubv4.NewDateTime(githubv4.DateTime{Time: toTime})
	}

	if d.EqualsQuals["max_repositories"] != nil {
		maxRepositories = int(d.EqualsQuals["max_repositories"].GetInt64Value())
	}

	if maxRepositories <= 0 {
		return nil, fmt.Errorf("invalid value for 'max_repositories' must be greater than 0")
	}

	var query struct {
		RateLimit models.RateLimit
		User      struct {
			ContributionsCollection models.ContributionsCollection `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":           githubv4.String(login),
		"from":            (*githubv4.DateTime)(nil),
		"to":              (*githubv4.DateTime)(nil),
		"maxRepositories": githubv4.Int(maxRepositories),
	}

	if fromDate != nil {
		variables["from"] = fromDate
	}
	if toDate != nil {
		variables["to"] = toDate
	}

	appendContributionColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_user_contribution_stats", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_user_contribution_stats", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, query.User.ContributionsCollection)

	return nil, nil
}
