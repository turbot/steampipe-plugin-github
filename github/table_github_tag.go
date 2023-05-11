package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTag() *plugin.Table {
	return &plugin.Table{
		Name:        "github_tag",
		Description: "Tags for commits in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTagList,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tag."},
			{Name: "tagger_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("TaggerDate").NullIfZero(), Description: "Date the tag was created."},
			{Name: "tagger_name", Type: proto.ColumnType_STRING, Description: "Name of user whom created the tag."},
			{Name: "tagger_login", Type: proto.ColumnType_STRING, Description: "Login of user whom created the tag."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Message associated with the tag."},
			{Name: "commit", Type: proto.ColumnType_JSON, Description: "Commit the tag is associated with."},
		},
	}
}

func tableGitHubTagList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Refs struct {
				TotalCount int
				PageInfo   models.PageInfo
				Nodes      []models.TagWithCommits
			} `graphql:"refs(refPrefix: \"refs/tags/\", first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_tag", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_tag", "api_error", err)
			return nil, err
		}

		for _, tag := range query.Repository.Refs.Nodes {
			d.StreamListItem(ctx, mapTagRow(&tag))

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(query.Repository.Refs.PageInfo.EndCursor)
	}

	return nil, nil
}

// tagRow is a struct to flatten returned information.
type tagRow struct {
	Name        string
	TaggerDate  time.Time
	TaggerName  string
	TaggerLogin string
	Message     string
	Commit      models.Commit
}

// mapTagRow is required as commit information may reside at upper target level or embedded into the tags target level.
func mapTagRow(tag *models.TagWithCommits) tagRow {
	row := tagRow{
		Name:        tag.Name,
		TaggerName:  tag.Target.Tag.Tagger.Name,
		TaggerDate:  tag.Target.Tag.Tagger.Date,
		TaggerLogin: tag.Target.Tag.Tagger.User.Login,
		Message:     tag.Target.Tag.Message,
	}

	if tag.Target.Commit.Sha != "" {
		row.Commit = tag.Target.Commit
	} else {
		row.Commit = tag.Target.Tag.Target.Commit
	}

	return row
}
