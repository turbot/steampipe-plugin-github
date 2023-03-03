package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubStargazer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_stargazer",
		Description: "Stargazers are users who have starred the repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubStargazerList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the stargazer."},
			{Name: "starred_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("StarredAt").Transform(convertTimestamp), Description: "Time when the stargazer was created."},
			{Name: "user_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Login"), Description: "The login name of the user who starred the repository."},
			{Name: "user_company", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Company"), Description: "The company of the user who starred the repository."},
			{Name: "user_location", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Location"), Description: "The location of the user who starred the repository."},
			// No extra useful data over login - {Name: "user", Type: proto.ColumnType_JSON, Transform: transform.FromField("User"), Description: "Details of the user who starred the repository."},
		},
	}
}

var stargazersQuery struct {
	Repository struct {
		Stargazers struct {
			Edges []struct {
				Node struct {
					Login    githubv4.String
					Company  githubv4.String
					Location githubv4.String
				}
				StarredAt githubv4.DateTime
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage githubv4.Boolean
			}
		} `graphql:"stargazers(first: $stargazersPageSize, after: $stargazersCursor)"`
	} `graphql:"repository(name: $repositoryName, owner: $repositoryOwner)"`
}

//// LIST FUNCTION

func tableGitHubStargazerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)
	quals := d.KeyColumnQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	pageSize := 100
	limit := d.QueryContext.Limit

	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}
	variables := map[string]interface{}{
		"repositoryOwner":    githubv4.String(owner),
		"repositoryName":     githubv4.String(repo),
		"stargazersPageSize": githubv4.Int(pageSize),
		"stargazersCursor":   (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(ctx, &stargazersQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_stargazer", "api_error", err)
			// if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			// 	return nil, nil
			// }
			return nil, err
		}
		for _, member := range stargazersQuery.Repository.Stargazers.Edges {
			d.StreamListItem(ctx, member)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !stargazersQuery.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}
		variables["stargazersCursor"] = githubv4.NewString(stargazersQuery.Repository.Stargazers.PageInfo.EndCursor)
	}

	return nil, nil
}
