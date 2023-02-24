package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubMyRepositoryGraphql() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository_graphql",
		Description: "GitHub Repositories that you are associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyRepositoryGraphqlList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "visibility", Require: plugin.Optional},
			},
		},
		Columns: gitHubRepositoryGraphqlColumns(),
	}
}

func gitHubRepositoryGraphqlColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "raw", Type: proto.ColumnType_JSON, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromValue()},
		{Name: "full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromField("NameWithOwner")},
		{Name: "language", Type: proto.ColumnType_STRING, Description: "The repository language (JavaScript, Go, etc)", Transform: transform.FromField("PrimaryLanguage.Name")},
		{Name: "private", Type: proto.ColumnType_BOOL, Description: "If true, the repo is private, otherwise it is public.", Transform: transform.FromField("IsPrivate")},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the repo.", Transform: transform.FromField("Url")},
		// Other columns
		{Name: "allow_merge_commit", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows merge commits.", Transform: transform.FromField("MergeCommitAllowed")},
		{Name: "allow_rebase_merge", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows rebase merges.", Transform: transform.FromField("RebaseMergeAllowed")},
		{Name: "allow_squash_merge", Type: proto.ColumnType_BOOL, Description: "If true, the repository allows squash merges.", Transform: transform.FromField("SquashMergeAllowed")},
		{Name: "archived", Type: proto.ColumnType_BOOL, Description: "If true, the repository is archived and read-only.", Transform: transform.FromField("IsArchived")},
		// {Name: "clone_url", Type: proto.ColumnType_STRING, Description: "URL that can be provided to git clone to clone the repository via HTTPS."},
		{Name: "code_of_conduct_key", Type: proto.ColumnType_STRING, Description: "Unique key for code of conduct for the repository.", Transform: transform.FromField("CodeOfConduct.Key")},
		{Name: "code_of_conduct_name", Type: proto.ColumnType_STRING, Description: "Name of the Code of Conduct for the repository.", Transform: transform.FromField("CodeOfConduct.Name")},
		{Name: "code_of_conduct_url", Type: proto.ColumnType_STRING, Description: "URL of the Code of Conduct for the repository.", Transform: transform.FromField("CodeOfConduct.Url")},
		{Name: "collaborators", Type: proto.ColumnType_JSON, Description: "An array of users (teams and outside collaborators) who have access to the repository, including their permissions.", Hydrate: getRepositoryCollaborators},
		// {Name: "collaborator_logins", Type: proto.ColumnType_JSON, Description: "An array of logins for users (inside and outside collaborators) who have access to the repository.", Transform: transform.FromValue().Transform(filterUserLogins), Hydrate: tableGitHubRepositoryCollaboratorsGetAll},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the repository was created.", Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
		{Name: "default_branch", Type: proto.ColumnType_STRING, Description: "The name of the deafult branch. The default branch is the base branch for pull requests and code commits.", Transform: transform.FromField("DefaultBranchRef.Name")},
		{Name: "delete_branch_on_merge", Type: proto.ColumnType_BOOL, Description: "If enabled, branches are automatically deleted whe a PR is merged."},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The repository description."},
		{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "If true, the repository is disabled.", Transform: transform.FromField("IsDisabled")},
		{Name: "fork", Type: proto.ColumnType_BOOL, Description: "If true, this repository is a fork of another repository.", Transform: transform.FromField("IsFork")},
		{Name: "forks_count", Type: proto.ColumnType_INT, Description: "The number of repositories that have forked this repository.", Transform: transform.FromField("ForkCount")},
		// {Name: "git_url", Type: proto.ColumnType_STRING, Description: "The git url to clone this repo via the git protocol."},
		// {Name: "has_downloads", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Downloads feature is enabled on the repository."},
		{Name: "has_issues", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Issues feature is enabled on the repository.", Transform: transform.FromField("HasIssuesEnabled")},
		// {Name: "has_pages", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Pages feature is enabled on the repository."},
		{Name: "has_projects", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Projects feature is enabled on the repository.", Transform: transform.FromField("HasProjectsEnabled")},
		{Name: "has_wiki", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Wiki feature is enabled on the repository.", Transform: transform.FromField("HasWikiEnabled")},
		{Name: "homepage", Type: proto.ColumnType_STRING, Description: "The URL of a page describing the project.", Transform: transform.FromField("HomepageUrl")},
		// {Name: "hooks", Type: proto.ColumnType_JSON, Description: "The API Hooks URL.", Hydrate: repositoryHooksGet, Transform: transform.FromValue()},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the repository."},
		{Name: "is_template", Type: proto.ColumnType_BOOL, Description: "If true, the repository is a template repository."},
		{Name: "license_key", Type: proto.ColumnType_STRING, Description: "The key of the license associated with the repository.", Transform: transform.FromField("LicenseInfo.Key")},
		{Name: "license_name", Type: proto.ColumnType_STRING, Description: "The name of the license associated with the repository.", Transform: transform.FromField("LicenseInfo.Name")},
		// {Name: "license_node_id", Type: proto.ColumnType_STRING, Description: "The node id of the license associated with the repository.", Transform: transform.FromField("License.NodeID")},
		{Name: "license_spdx_id", Type: proto.ColumnType_STRING, Description: "The Software Package Data Exchange (SPDX) id of the license associated with the repository.", Transform: transform.FromField("LicenseInfo.SpdxId")},
		{Name: "license_url", Type: proto.ColumnType_STRING, Description: "The url of the license associated with the repository.", Transform: transform.FromField("LicenseInfo.Url")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the repository."},
		// {Name: "network_count", Type: proto.ColumnType_INT, Description: "The number of member repositories in the network.", Hydrate: tableGitHubRepositoryGet},
		// {Name: "node_id", Type: proto.ColumnType_STRING, Description: "The Node ID of the repository."},
		{Name: "open_issues_count", Type: proto.ColumnType_INT, Description: "The number of open issues for the repository.", Transform: transform.FromField("Issues.TotalCount")},
		// Only load relevant fields from the owner
		{Name: "owner_id", Type: proto.ColumnType_STRING, Description: "The user id (number) of the repository owner.", Transform: transform.FromField("Owner.Id")},
		{Name: "owner_login", Type: proto.ColumnType_STRING, Description: "The user login name of the repository owner.", Transform: transform.FromField("Owner.Login")},
		{Name: "owner_type", Type: proto.ColumnType_STRING, Description: "The type of the repository owner (User or Organization).", Transform: transform.FromField("Owner.__typename")},
		// {Name: "outside_collaborators", Type: proto.ColumnType_JSON, Description: "An array of outside collaborators who have access to the repository, including their permissions.", Transform: transform.FromValue(), Hydrate: tableGitHubRepositoryCollaboratorsGetOutside},
		// {Name: "outside_collaborator_logins", Type: proto.ColumnType_JSON, Description: "An array of logins for outside collaborators who have access to the repository.", Transform: transform.FromValue().Transform(filterUserLogins), Hydrate: tableGitHubRepositoryCollaboratorsGetOutside},
		{Name: "pushed_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of the last push to the repository.", Transform: transform.FromField("PushedAt").Transform(convertTimestamp)},
		{Name: "size", Type: proto.ColumnType_INT, Description: "The size of the whole repository (including history), in kilobytes.", Transform: transform.FromField("DiskUsage")},
		{Name: "ssh_url", Type: proto.ColumnType_STRING, Description: "The url to clone this repo via ssh."},
		{Name: "stargazers_count", Type: proto.ColumnType_INT, Description: "The number of users who have 'starred' the repository."},
		// {Name: "subscribers_count", Type: proto.ColumnType_INT, Description: "The number of users who have subscribed to the repository.", Hydrate: tableGitHubRepositoryGet},
		{Name: "template_repository", Type: proto.ColumnType_STRING, Description: "The template repository used to create this resource.", Transform: transform.FromField("TemplateRepository.NameWithOwner")},
		// {Name: "topics", Type: proto.ColumnType_JSON, Description: "The topics (similar to tags or labels) associated with the repository."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the repository was last updated.", Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
		// {Name: "url", Type: proto.ColumnType_STRING, Description: "The url to clone this repo via https."},
		{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the repository (public or private)"},
		{Name: "watchers_count", Type: proto.ColumnType_INT, Description: "The number of users who have watched the repository.", Transform: transform.FromField("Watchers.TotalCount")},
	}
}

//// LIST FUNCTION

/* func tableGitHubMyRepositoryGraphqlList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	// Additional filters
	if d.KeyColumnQuals["visibility"] != nil {
		opt.Visibility = d.KeyColumnQuals["visibility"].GetStringValue()
	} else {
		// Will cause a 422 error if 'type' used in the same request as visibility or
		// affiliation.
		opt.Type = "all"
	}
	type ListPageResponse struct {
		repo []*github.Repository
		resp *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		return ListPageResponse{
			repo: repos,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repos := listResponse.repo
		resp := listResponse.resp

		for _, i := range repos {
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

		opt.Page = resp.NextPage
	}

	return nil, nil
} */

var MyRepositoryQuery struct {
	Viewer struct {
		Login        githubv4.String
		Repositories struct {
			Edges    []repositories
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: $repositoriesPageSize, after: $repositoriesCursor, ownerAffiliations: $ownerAffiliations, affiliations: $affiliations)"`
	}
}

type repositories struct {
	Node Node
}

type Node struct {
	NameWithOwner   githubv4.String
	PrimaryLanguage struct {
		Name githubv4.String
	}
	IsPrivate          bool
	Url                githubv4.String
	MergeCommitAllowed bool
	RebaseMergeAllowed bool
	SquashMergeAllowed bool
	IsArchived         bool
	CodeOfConduct      struct {
		Id   githubv4.String
		Key  githubv4.String
		Name githubv4.String
		Url  githubv4.String
	}
	CreatedAt        githubv4.GitTimestamp
	DefaultBranchRef struct {
		Name githubv4.String
	}
	DeleteBranchOnMerge bool
	Description         githubv4.String
	IsDisabled          bool
	IsFork              bool
	ForkCount           githubv4.Int
	HasIssuesEnabled    bool
	HasProjectsEnabled  bool
	HasWikiEnabled      bool
	HomepageUrl         githubv4.String
	Id                  githubv4.String
	IsTemplate          bool
	LicenseInfo         struct {
		Id             githubv4.String
		Key            githubv4.String
		SpdxId         githubv4.String
		Url            githubv4.String
		PseudoLicense  githubv4.Boolean
		Nickname       githubv4.String
		Name           githubv4.String
		Implementation githubv4.String
		Hidden         githubv4.Boolean
		Featured       githubv4.Boolean
		Description    githubv4.String
		Body           githubv4.String
	}
	Name   githubv4.String
	Issues struct {
		TotalCount githubv4.Int
	} `graphql:"issues(states: OPEN)"`
	Owner struct {
		Id    githubv4.String
		Login githubv4.String
	}
	PushedAt           githubv4.String
	DiskUsage          githubv4.Int
	SshUrl             githubv4.String
	StargazerCount     githubv4.Int
	TemplateRepository struct {
		NameWithOwner githubv4.String
	}
	UpdatedAt  githubv4.GitTimestamp
	Visibility githubv4.String
	Watchers   struct {
		TotalCount githubv4.Int
	}
}

// var Collaborators struct {
// 	Edges    []collaborators
// 	PageInfo struct {
// 		EndCursor   githubv4.String
// 		HasNextPage bool
// 	}
// } `graphql:"collaborators(affiliation: ALL)"`

var MyRepositoryCollaboratorsQuery struct {
	Repository struct {
		Collaborators struct {
			TotalCount githubv4.Int
			Edges      []collaborators
		} `graphql:"collaborators(affiliation: ALL)"`
	} `graphql:"repository(name: $repositoryName, owner: $repositoryOwner)"`
}

type collaborators struct {
	Node struct {
		Id          githubv4.String
		Url         githubv4.String
		Login       githubv4.String
		AvatarUrl   githubv4.String `graphql:"avatarUrl(size: 10)"`
		IsSiteAdmin bool
	}
	Permission        githubv4.String
	PermissionSources struct {
		Permission githubv4.String
		Source     githubv4.String
		// __typename githubv4.String
	}
}

func tableGitHubMyRepositoryGraphqlList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	pageSize := 100

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"repositoriesPageSize": githubv4.Int(pageSize),
		"repositoriesCursor":   (*githubv4.String)(nil), // Null after argument to get first page.
		"ownerAffiliations":    []githubv4.RepositoryAffiliation{githubv4.RepositoryAffiliationOrganizationMember},
		"affiliations":         []githubv4.RepositoryAffiliation{githubv4.RepositoryAffiliationOrganizationMember},
		// "collaboratorsPageSize": githubv4.Int(pageSize),
	}

	for {
		err := client.Query(ctx, &MyRepositoryQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_member", "api_error", err)
			// if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			// 	return nil, nil
			// }
			return nil, err
		}

		for _, repo := range MyRepositoryQuery.Viewer.Repositories.Edges {
			d.StreamListItem(ctx, repo.Node)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !MyRepositoryQuery.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["repositoriesCursor"] = githubv4.NewString(MyRepositoryQuery.Viewer.Repositories.PageInfo.EndCursor)
	}

	return nil, nil
}

func getRepositoryCollaborators(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)
	repo := h.Item.(Node)

	plugin.Logger(ctx).Trace("------------------------>>>>>>>>>>>>>>>>>>>", repo.Name, repo.Owner.Login)
	variables := map[string]interface{}{
		"repositoryName":  *githubv4.NewString(repo.Name),
		"repositoryOwner": *githubv4.NewString(repo.Owner.Login),
	}

	err := client.Query(ctx, &MyRepositoryCollaboratorsQuery, variables)
	if err != nil {
		plugin.Logger(ctx).Error("github_organization_member", "api_error", err)
		// if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
		// 	return nil, nil
		// }
		return nil, err
	}
	return MyRepositoryCollaboratorsQuery.Repository.Collaborators.Edges, nil
}
