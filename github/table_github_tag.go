package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_tag",
		Description: "Tags for commits in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTagList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tag."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Target.Tag.Tagger.Date").NullIfZero(), Description: "Date the tag was created."},
			{Name: "tagger_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Tagger.Name"), Description: "Name of user whom created the tag."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.Oid"), Description: "Commit SHA the tag refers to."},
			{Name: "commit_short_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.AbbreviatedOid"), Description: "Commit short SHA the tag refers to."},
			{Name: "commit_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Target.Tag.Target.Commit.CommittedDate").NullIfZero(), Description: "Date the commit referenced by the tag was made."},
			{Name: "commit_message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.Message"), Description: "Commit message of the commit the tag is referencing."},
			{Name: "commit_author", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.Committer.Name"), Description: "Name of the author for the commit the tag is referencing."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.CommitUrl"), Description: "Commit URL the tag refers to."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.ZipballUrl"), Description: "URL to download a zip file of the code for this tag."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Target.Tag.Target.Commit.TarballUrl"), Description: "URL to download a tar file of the code for this tag."},
		},
	}
}

type tagDetail struct {
	Name   string
	Target struct {
		Tag struct {
			Tagger struct {
				Name string
				Date time.Time
			}
			CommitUrl string
			Message   string
			Target    struct {
				Commit struct {
					Oid            string
					AbbreviatedOid string
					CommittedDate  time.Time
					Message        string
					TarballUrl     string
					ZipballUrl     string
					Committer      struct {
						Name string
					}
				} `graphql:"... on Commit"`
			}
		} `graphql:"... on Tag"`
	}
}

var tagQuery struct {
	Repository struct {
		Refs struct {
			TotalCount int
			PageInfo   struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
			Nodes []tagDetail
		} `graphql:"refs(refPrefix: \"refs/tags/\", first: $tagsPageSize, after: $tagsCursor)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

func tableGitHubTagList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := 100
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"tagsPageSize":    githubv4.Int(pageSize),
		"tagsCursor":      (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &tagQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_tag", "api_error", err)
			return nil, err
		}

		for _, tag := range tagQuery.Repository.Refs.Nodes {
			d.StreamListItem(ctx, tag)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !tagQuery.Repository.Refs.PageInfo.HasNextPage {
			break
		}

		variables["tagsCursor"] = githubv4.NewString(tagQuery.Repository.Refs.PageInfo.EndCursor)
	}

	return nil, nil
}
