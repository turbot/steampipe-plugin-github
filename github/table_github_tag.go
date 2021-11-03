package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_tag",
		Description: "Tags for commits in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubTagList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tag."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.SHA"), Description: "Commit SHA the tag refers to."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.URL"), Description: "Commit URL the tag refers to."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING, Description: "URL to download a zip file of the code for this tag."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING, Description: "URL to download a tar file of the code for this tag."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubTagList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	type ListPageResponse struct {
		tags []*github.RepositoryTag
		resp *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opts)
		return ListPageResponse{
			tags: tags,
			resp: resp,
		}, err
	}

	for {

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		tags := listResponse.tags
		resp := listResponse.resp

		for _, i := range tags {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return nil, nil
}
