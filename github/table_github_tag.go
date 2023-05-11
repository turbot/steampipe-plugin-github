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
			{Name: "tag_message", Type: proto.ColumnType_STRING, Description: "Message associated with the tag."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Description: "Commit SHA the branch refers to."},
			{Name: "commit_short_sha", Type: proto.ColumnType_STRING, Description: "Commit short SHA the branch refers to."},
			{Name: "commit_authored_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date commit was authored."},
			{Name: "commit_author_name", Type: proto.ColumnType_STRING, Description: "Commit authors display name."},
			{Name: "commit_author_login", Type: proto.ColumnType_STRING, Description: "Commit authors login."},
			{Name: "commit_committed_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date commit was committed."},
			{Name: "commit_committer_name", Type: proto.ColumnType_STRING, Description: "Commit committers display name."},
			{Name: "commit_committer_login", Type: proto.ColumnType_STRING, Description: "Commit committers login."},
			{Name: "commit_message", Type: proto.ColumnType_STRING, Description: "Commit message."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Description: "Commit URL the branch refers to."},
			{Name: "commit_additions", Type: proto.ColumnType_INT, Description: "Number of additions in the commit."},
			{Name: "commit_deletions", Type: proto.ColumnType_INT, Description: "Number of deletions in the commit."},
			{Name: "commit_changed_files", Type: proto.ColumnType_INT, Transform: transform.FromField("CommitChangedFiles").NullIfZero(), Description: "Number of files changed in the commit if available (null if not available)."},
			{Name: "commit_authored_by_committer", Type: proto.ColumnType_BOOL, Description: "If true, the commits committer and author are the same."},
			{Name: "commit_committed_via_web", Type: proto.ColumnType_BOOL, Description: "If true, the commit was from the GitHub web app."},
			{Name: "commit_signature_is_valid", Type: proto.ColumnType_BOOL, Description: "If true, commit was signed by a valid signature."},
			{Name: "commit_signature_email", Type: proto.ColumnType_STRING, Description: "Email associated with the commit signature."},
			{Name: "commit_signature_login", Type: proto.ColumnType_STRING, Description: "Login associated with the commit signature."},
			{Name: "commit_tarball_url", Type: proto.ColumnType_STRING, Description: "URL to download a tar file of the code for this commit."},
			{Name: "commit_zipball_url", Type: proto.ColumnType_STRING, Description: "URL to download a zip file of the code for this commit."},
			{Name: "commit_tree_url", Type: proto.ColumnType_STRING, Description: "URL for the tree of this commit."},
			{Name: "commit_status", Type: proto.ColumnType_STRING, Description: "Status of the commit (ERROR, EXPECTED, FAILURE, PENDING, SUCCESS)."},
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
	Name                      string
	TaggerDate                time.Time
	TaggerName                string
	TaggerLogin               string
	TagMessage                string
	CommitSha                 string
	CommitShortSha            string
	CommitAuthoredDate        time.Time
	CommitAuthorName          string
	CommitAuthorLogin         string
	CommitCommittedDate       time.Time
	CommitCommitterName       string
	CommitCommitterLogin      string
	CommitMessage             string
	CommitURL                 string
	CommitAdditions           int
	CommitDeletions           int
	CommitChangedFiles        int
	CommitAuthoredByCommitter bool
	CommitCommittedViaWeb     bool
	CommitSignatureIsValid    bool
	CommitSignatureEmail      string
	CommitSignatureLogin      string
	CommitTreeURL             string
	CommitZipballURL          string
	CommitTarballURL          string
	CommitStatus              string
}

// mapTagRow is required as commit information may reside at upper target level or embedded into the tags target level.
func mapTagRow(tag *models.TagWithCommits) tagRow {
	row := tagRow{
		Name:        tag.Name,
		TaggerName:  tag.Target.Tag.Tagger.Name,
		TaggerDate:  tag.Target.Tag.Tagger.Date,
		TaggerLogin: tag.Target.Tag.Tagger.User.Login,
		TagMessage:  tag.Target.Tag.Message,
	}

	if tag.Target.Commit.Sha != "" {
		mapTagRowCommitParser(&row, &tag.Target.Commit)
	} else {
		mapTagRowCommitParser(&row, &tag.Target.Tag.Target.Commit)
	}

	return row
}

func mapTagRowCommitParser(row *tagRow, commit *models.Commit) {
	row.CommitSha = commit.Sha
	row.CommitShortSha = commit.ShortSha
	row.CommitAuthoredDate = commit.AuthoredDate
	row.CommitAuthorName = commit.Author.Name
	row.CommitAuthorLogin = commit.Author.User.Login
	row.CommitCommittedDate = commit.CommittedDate
	row.CommitCommitterName = commit.Committer.Name
	row.CommitCommitterLogin = commit.Committer.User.Login
	row.CommitMessage = commit.Message
	row.CommitURL = commit.Url
	row.CommitAdditions = commit.Additions
	row.CommitDeletions = commit.Deletions
	row.CommitChangedFiles = commit.ChangedFiles
	row.CommitAuthoredByCommitter = commit.AuthoredByCommitter
	row.CommitCommittedViaWeb = commit.CommittedViaWeb
	row.CommitSignatureIsValid = commit.Signature.IsValid
	row.CommitSignatureEmail = commit.Signature.Email
	row.CommitSignatureLogin = commit.Signature.Signer.Login
	row.CommitTreeURL = commit.TreeUrl
	row.CommitZipballURL = commit.ZipballUrl
	row.CommitTarballURL = commit.TarballUrl
	row.CommitStatus = commit.Status.State
}
