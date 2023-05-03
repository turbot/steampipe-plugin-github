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
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Date").NullIfZero(), Description: "Date the tag was created."},
			{Name: "tagger_name", Type: proto.ColumnType_STRING, Description: "Name of user whom created the tag."},
			{Name: "tag_message", Type: proto.ColumnType_STRING, Description: "Message associated with the tag."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Description: "Commit SHA the tag refers to."},
			{Name: "commit_short_sha", Type: proto.ColumnType_STRING, Description: "Commit short SHA the tag refers to."},
			{Name: "commit_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CommitDate").NullIfZero(), Description: "Date the commit referenced by the tag was made."},
			{Name: "commit_message", Type: proto.ColumnType_STRING, Description: "Commit message of the commit the tag is referencing."},
			{Name: "commit_author", Type: proto.ColumnType_STRING, Description: "Name of the author for the commit the tag is referencing."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("CommitUrl"), Description: "Commit URL the tag refers to."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZipballUrl"), Description: "URL to download a zip file of the code for this tag."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TarballUrl"), Description: "URL to download a tar file of the code for this tag."},
		},
	}
}

type tagRow struct {
	Name           string
	Date           time.Time
	TaggerName     string
	TagMessage     string
	CommitSha      string
	CommitShortSha string
	CommitDate     time.Time
	CommitMessage  string
	CommitAuthor   string
	CommitUrl      string
	ZipballUrl     string
	TarballUrl     string
}

type tagDetail struct {
	Name   string
	Target struct {
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
			CommitUrl string
		} `graphql:"... on Commit"`
		Tag struct {
			Tagger struct {
				Name string
				Date time.Time
			}
			Message string
			Target  struct {
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
					CommitUrl string
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
			d.StreamListItem(ctx, generateTagRow(tag))

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

// generateTagRow is used to parse the tagDetail into a common return type tagRow as data can be in different places depending on how tag was built.
func generateTagRow(tag tagDetail) tagRow {
	row := tagRow{
		Name:       tag.Name,
		TaggerName: tag.Target.Tag.Tagger.Name,
		Date:       tag.Target.Tag.Tagger.Date,
		TagMessage: tag.Target.Tag.Message,
	}

	if tag.Target.Commit.Oid != "" {
		row.CommitSha = tag.Target.Commit.Oid
		row.CommitShortSha = tag.Target.Commit.AbbreviatedOid
		row.CommitDate = tag.Target.Commit.CommittedDate
		row.CommitAuthor = tag.Target.Commit.Committer.Name
		row.CommitMessage = tag.Target.Commit.Message
		row.CommitUrl = tag.Target.Commit.CommitUrl
		row.TarballUrl = tag.Target.Commit.TarballUrl
		row.ZipballUrl = tag.Target.Commit.ZipballUrl
	} else {
		row.CommitSha = tag.Target.Tag.Target.Commit.Oid
		row.CommitShortSha = tag.Target.Tag.Target.Commit.AbbreviatedOid
		row.CommitDate = tag.Target.Tag.Target.Commit.CommittedDate
		row.CommitAuthor = tag.Target.Tag.Target.Commit.Committer.Name
		row.CommitMessage = tag.Target.Tag.Target.Commit.Message
		row.CommitUrl = tag.Target.Tag.Target.Commit.CommitUrl
		row.TarballUrl = tag.Target.Tag.Target.Commit.TarballUrl
		row.ZipballUrl = tag.Target.Tag.Target.Commit.ZipballUrl
	}
	return row
}
