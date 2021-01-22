package github

import (
	"context"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository",
		Description: "Github Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubRepositoryList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"owner_login", "name"}),
			Hydrate:    tableGitHubRepositoryGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name."},
			{Name: "language", Type: proto.ColumnType_STRING, Description: "The repository language (JavaScript, Go, etc)"},
			{Name: "private", Type: proto.ColumnType_BOOL, Description: "If true, the repo is private, otherwise it is public."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the repo."},

			{Name: "allow_merge_commit", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows merge commits.", Hydrate: tableGitHubRepositoryGet},
			{Name: "allow_rebase_merge", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows rebase merges.", Hydrate: tableGitHubRepositoryGet},
			{Name: "allow_squash_merge", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows squash merges.", Hydrate: tableGitHubRepositoryGet},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows rebase merges."},
			{Name: "clone_url", Type: proto.ColumnType_STRING, Description: "URL that can be provided to git clone to clone the repository via HTTPS."},
			{Name: "code_of_conduct_key", Type: proto.ColumnType_STRING, Description: "Unique key for code of conduct for the repository.", Transform: transform.FromField("CodeOfConduct.Key"), Hydrate: tableGitHubRepositoryGet},
			{Name: "code_of_conduct_name", Type: proto.ColumnType_STRING, Description: "Name of the Code of Conduct for the repository.", Transform: transform.FromField("CodeOfConduct.Name"), Hydrate: tableGitHubRepositoryGet},
			{Name: "code_of_conduct_url", Type: proto.ColumnType_STRING, Description: "URL of the Code of Conduct for the repository.", Transform: transform.FromField("CodeOfConduct.URL"), Hydrate: tableGitHubRepositoryGet},
			{Name: "collaborators", Type: proto.ColumnType_JSON, Description: "An array of users who have access to the repository, including their permissions.", Transform: transform.FromValue(), Hydrate: tableGitHubRepositoryCollaboratorsGet},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the repository was created.", Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
			{Name: "default_branch", Type: proto.ColumnType_STRING, Description: "The name of the deafult branch. The default branch is the base branch for pull requests and code commits."},
			{Name: "delete_branch_on_merge", Type: proto.ColumnType_BOOL, Description: "If enabled, branches are automatically deleted whe a PR is merged.", Hydrate: tableGitHubRepositoryGet},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The repository description."},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "If true, the repository is disabled."},
			{Name: "fork", Type: proto.ColumnType_BOOL, Description: "If true, this repository is a fork of another repository."},
			{Name: "forks_count", Type: proto.ColumnType_INT, Description: "The number of repositories that have forked this repository."},
			{Name: "git_url", Type: proto.ColumnType_STRING, Description: "The git url to clone this repo via the git protocol."},
			{Name: "has_downloads", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Downloads feature is enabled on the repository."},
			{Name: "has_issues", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Issues feature is enabled on the repository."},
			{Name: "has_pages", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Pages feature is enabled on the repository."},
			{Name: "has_projects", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Projects feature is enabled on the repository."},
			{Name: "has_wiki", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Wiki feature is enabled on the repository."},
			{Name: "homepage", Type: proto.ColumnType_STRING, Description: "The URL of a page describing the project."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the repository."},
			{Name: "is_template", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubRepositoryGet, Description: "If true, the repository is a template repository."},
			{Name: "license_key", Type: proto.ColumnType_STRING, Description: "The key of the license associated with the repository.", Transform: transform.FromField("License.Key")},
			{Name: "license_name", Type: proto.ColumnType_STRING, Description: "The name of the license associated with the repository.", Transform: transform.FromField("License.Name")},
			{Name: "license_node_id", Type: proto.ColumnType_STRING, Description: "The node id of the license associated with the repository.", Transform: transform.FromField("License.NodeID")},
			{Name: "license_spdx_id", Type: proto.ColumnType_STRING, Description: "The Software Package Data Exchange (SPDX) id of the license associated with the repository.", Transform: transform.FromField("License.SPDXID")},
			{Name: "license_url", Type: proto.ColumnType_STRING, Description: "The url of the license associated with the repository.", Transform: transform.FromField("License.URL")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the repository."},
			{Name: "network_count", Type: proto.ColumnType_INT, Description: "The number of member repositories in the network.", Hydrate: tableGitHubRepositoryGet},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The Node ID of the repository."},
			{Name: "open_issues_count", Type: proto.ColumnType_INT, Description: "The number of open issues for the repository."},
			// Only load relevant fields from the owner
			{Name: "owner_id", Type: proto.ColumnType_INT, Description: "The user id (number) of the repository owner.", Transform: transform.FromField("Owner.ID")},
			{Name: "owner_login", Type: proto.ColumnType_STRING, Description: "The user login name of the repository owner.", Transform: transform.FromField("Owner.Login")},
			{Name: "owner_type", Type: proto.ColumnType_STRING, Description: "The type of the repository owner (User or Organization).", Transform: transform.FromField("Owner.Type")},
			{Name: "pushed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of the last push to the repository.", Transform: transform.FromField("PushedAt").Transform(convertTimestamp)},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The size of the whole repository (including history), in kilobytes."},
			{Name: "ssh_url", Type: proto.ColumnType_STRING, Description: "The url to clone this repo via ssh."},
			{Name: "stargazers_count", Type: proto.ColumnType_INT, Description: "The number of users who have 'starred' the repository."},
			{Name: "subscribers_count", Type: proto.ColumnType_INT, Description: "The number of users who have subscribed to the repository.", Hydrate: tableGitHubRepositoryGet},
			{Name: "template_repository", Type: proto.ColumnType_STRING, Description: "The template repository used to create this resource."},
			{Name: "topics", Type: proto.ColumnType_JSON, Description: "The topics (similar to tags or labels) associated with the repository."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the repository was last updated.", Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The url to clone this repo via https."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the repository (public or private)", Hydrate: tableGitHubRepositoryGet},
			{Name: "watchers_count", Type: proto.ColumnType_INT, Description: "The number of users who have watched the repository."},
		},
	}
}

type gitHubRepositoryCollaborator struct {
	Login string
}

//// list ////

func tableGitHubRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d.ConnectionManager)

	opt := &github.RepositoryListOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 100}}

	for {

		var repos []*github.Repository
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			repos, resp, err = client.Repositories.List(ctx, "", opt)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range repos {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

//// hydrate functions ////

func tableGitHubRepositoryGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var owner, repoName string
	if h.Item != nil {
		repo := h.Item.(*github.Repository)
		owner = *repo.Owner.Login
		repoName = *repo.Name
	} else {
		owner = d.KeyColumnQuals["owner_login"].GetStringValue()
		repoName = d.KeyColumnQuals["name"].GetStringValue()
	}

	client := connect(ctx, d.ConnectionManager)

	var detail *github.Repository
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		detail, resp, err = client.Repositories.Get(ctx, owner, repoName)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}

func tableGitHubRepositoryCollaboratorsGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	repo := h.Item.(*github.Repository)
	owner := *repo.Owner.Login
	repoName := *repo.Name

	client := connect(ctx, d.ConnectionManager)

	var repositoryCollaborators []*github.User

	opt := &github.ListCollaboratorsOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {

		var users []*github.User
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			users, resp, err = client.Repositories.ListCollaborators(ctx, owner, repoName, opt)
			logger.Trace("tableGitHubRepositoryCollaboratorsGet", "Users", users)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range users {
			repositoryCollaborators = append(repositoryCollaborators, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	logger.Trace("RepositoryCollaborators", repositoryCollaborators)

	return repositoryCollaborators, nil
}
