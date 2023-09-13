package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubCodeOwnerColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name."},
		// Other columns
		{Name: "line", Type: proto.ColumnType_INT, Description: "The rule's line number in the CODEOWNERS file.", Transform: transform.FromField("LineNumber")},
		{Name: "pattern", Type: proto.ColumnType_STRING, Description: "The pattern used to identify what code a team, or an individual is responsible for"},
		{Name: "users", Type: proto.ColumnType_JSON, Description: "Users responsible for code in the repo"},
		{Name: "teams", Type: proto.ColumnType_JSON, Description: "Teams responsible for code in the repo"},
		{Name: "pre_comments", Type: proto.ColumnType_JSON, Description: "Specifies the comments added above a key."},
		{Name: "line_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comment following the node and before empty lines."},
	}
}

func tableGitHubCodeOwner() *plugin.Table {
	return &plugin.Table{
		Name:        "github_code_owner",
		Description: "Individuals or teams that are responsible for code in a repository.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubCodeOwnerList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
		},
		Columns: gitHubCodeOwnerColumns(),
	}
}

type CodeOwnerRule struct {
	LineNumber  int
	Pattern     string
	Users       []string
	Teams       []string
	PreComments []string
	LineComment string
}

func tableGitHubCodeOwnerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableGitHubCodeOwnerList")
	repoFullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	type CodeOwnerRuleResponse struct {
		RepositoryFullName string
		LineNumber         int
		Pattern            string
		Users              []string
		Teams              []string
		PreComments        []string
		LineComment        string
	}

	getCodeOwners := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		var fileContent *github.RepositoryContent
		var err error

		client := connect(ctx, d)
		opt := &github.RepositoryContentGetOptions{}
		// If a repository has multiple CODEOWNERS files, GitHub uses the following
		// precedence: .github/CODEOWNERS, CODEOWNERS, docs/CODEOWNERS
		var paths = []string{".github/CODEOWNERS", "CODEOWNERS", "docs/CODEOWNERS"}
		for _, path := range paths {
			fileContent, _, _, err = client.Repositories.GetContents(ctx, owner, repoName, path, opt)
			// Stop on the first CODEOWNERS file found
			if err == nil {
				break
			}
			// HTTP 404 is the only tolerated HTTP error code, if it's different, it
			// means something is wrong with your rights or repository
			if err.(*github.ErrorResponse).Response.StatusCode != 404 {
				plugin.Logger(ctx).Error("github_code_owner.tableGitHubCodeOwnerList", "api_error", err, "path", path)
				return nil, err
			}
		}

		// No CODEOWNERS file
		if fileContent == nil {
			return []*CodeOwnerRule{}, nil
		}

		decodedContent, err := fileContent.GetContent()
		if err != nil {
			plugin.Logger(ctx).Error("github_code_owner.tableGitHubCodeOwnerList", "decode_error", err)
			return []*CodeOwnerRule{}, err
		}

		return decodeCodeOwnerFileContent(decodedContent), err
	}

	codeOwnersElements, err := plugin.RetryHydrate(ctx, d, h, getCodeOwners, retryConfig())
	if err != nil {
		plugin.Logger(ctx).Error("github_code_owner.tableGitHubCodeOwnerList", "retry_hydrate_error", err)
		return nil, err
	}

	for _, codeOwner := range codeOwnersElements.([]*CodeOwnerRule) {
		if codeOwner != nil {
			d.StreamListItem(ctx, CodeOwnerRuleResponse{
				RepositoryFullName: repoFullName,
				LineNumber:         codeOwner.LineNumber,
				Pattern:            codeOwner.Pattern,
				Users:              codeOwner.Users,
				Teams:              codeOwner.Teams,
				PreComments:        codeOwner.PreComments,
				LineComment:        codeOwner.LineComment,
			})
		}
	}
	return nil, nil
}

func decodeCodeOwnerFileContent(content string) []*CodeOwnerRule {
	var codeOwnerRules []*CodeOwnerRule

	var comments []string
	for i, line := range strings.Split(content, "\n") {
		lineNumber := i + 1
		// if line is empty, consider the codeblock end
		if len(line) == 0 {
			comments = []string{}
			continue
		}
		// code block with comments
		if strings.HasPrefix(line, "#") {
			comments = append(comments, line)
			continue
		}
		// code owners rule
		// if line is empty, consider the codeblock end
		rule := strings.SplitN(line, " ", 2)
		if len(rule) < 2 {
			comments = []string{}
			continue
		}

		var pattern, lineComment string
		pattern = rule[0]

		// line comment
		ownersAndComment := strings.SplitN(rule[1], "#", 2)
		if len(ownersAndComment) == 2 && len(ownersAndComment[1]) > 0 {
			lineComment = ownersAndComment[1]
		} else {
			ownersAndComment = []string{rule[1]}
		}

		// owners computing
		var users, teams []string
		for _, owner := range strings.Split(strings.TrimSpace(ownersAndComment[0]), " ") {
			if strings.Index(owner, "/") > 0 {
				teams = append(teams, owner)
			} else {
				users = append(users, owner)
			}
		}
		codeOwnerRules = append(codeOwnerRules, &CodeOwnerRule{LineNumber: lineNumber, Pattern: pattern, Users: users, Teams: teams, PreComments: comments, LineComment: lineComment})
	}
	return codeOwnerRules
}
