package github

import (
	"github.com/shurcooL/githubv4"
	"golang.org/x/exp/slices"
)

func appendDiscussionColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeDiscussionId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeDiscussionNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeDiscussionCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeDiscussionUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeDiscussionAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeDiscussionCategory"] = githubv4.Boolean(slices.Contains(cols, "category") || slices.Contains(cols, "category_name"))
	(*m)["includeDiscussionAnswer"] = githubv4.Boolean(slices.Contains(cols, "answer"))
	(*m)["includeDiscussionComments"] = githubv4.Boolean(slices.Contains(cols, "comments") || slices.Contains(cols, "comments_total_count"))
}
