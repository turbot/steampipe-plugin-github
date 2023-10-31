package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractCommitFromHydrateItem(h *plugin.HydrateData) (models.Commit, error) {
	if commit, ok := h.Item.(models.Commit); ok {
		return commit, nil
	} else {
		return models.Commit{}, fmt.Errorf("unable to parse hydrate item %v as an Commit", h.Item)
	}
}

func appendCommitColumnIncludes(m *map[string]interface{}, cols []string) {
	// For BasicCommit struct
	(*m)["includeCommitShortSha"] = githubv4.Boolean(slices.Contains(cols, "short_sha"))
	(*m)["includeCommitAuthoredDate"] = githubv4.Boolean(slices.Contains(cols, "authored_date"))
	(*m)["includeCommitAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeCommitCommittedDate"] = githubv4.Boolean(slices.Contains(cols, "committed_date"))
	(*m)["includeCommitCommitter"] = githubv4.Boolean(slices.Contains(cols, "committer") || slices.Contains(cols, "committer_login"))
	(*m)["includeCommitMessage"] = githubv4.Boolean(slices.Contains(cols, "message"))
	(*m)["includeCommitUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))

	// For Commit struct
	(*m)["includeCommitAdditions"] = githubv4.Boolean(slices.Contains(cols, "additions"))
	(*m)["includeCommitAuthoredByCommitter"] = githubv4.Boolean(slices.Contains(cols, "authored_by_committer"))
	(*m)["includeCommitChangedFiles"] = githubv4.Boolean(slices.Contains(cols, "changed_files"))
	(*m)["includeCommitCommittedViaWeb"] = githubv4.Boolean(slices.Contains(cols, "committed_via_web"))
	(*m)["includeCommitCommitUrl"] = githubv4.Boolean(slices.Contains(cols, "commit_url"))
	(*m)["includeCommitDeletions"] = githubv4.Boolean(slices.Contains(cols, "deletions"))
	(*m)["includeCommitSignature"] = githubv4.Boolean(slices.Contains(cols, "signature"))
	(*m)["includeCommitTarballUrl"] = githubv4.Boolean(slices.Contains(cols, "tarball_url"))
	(*m)["includeCommitTreeUrl"] = githubv4.Boolean(slices.Contains(cols, "tree_url"))
	(*m)["includeCommitCanSubscribe"] = githubv4.Boolean(slices.Contains(cols, "can_subscribe"))
	(*m)["includeCommitSubscription"] = githubv4.Boolean(slices.Contains(cols, "subscription"))
	(*m)["includeCommitZipballUrl"] = githubv4.Boolean(slices.Contains(cols, "zipball_url"))
	(*m)["includeCommitMessageHeadline"] = githubv4.Boolean(slices.Contains(cols, "message_headline"))
	(*m)["includeCommitStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeCommitNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
}

func commitHydrateAuthorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Author.User.Login, nil
}

func commitHydrateCommitterLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Committer.User.Login, nil
}

func commitHydrateShortSha(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.ShortSha, nil
}

func commitHydrateAuthoredDate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.AuthoredDate, nil
}

func commitHydrateAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Author, nil
}

func commitHydrateCommittedDate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.CommittedDate, nil
}

func commitHydrateCommitter(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Committer, nil
}

func commitHydrateMessage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Message, nil
}

func commitHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Url, nil
}

func commitHydrateAdditions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Additions, nil
}

func commitHydrateAuthoredByCommitter(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.AuthoredByCommitter, nil
}

func commitHydrateChangedFiles(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.ChangedFiles, nil
}

func commitHydrateCommittedViaWeb(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.CommittedViaWeb, nil
}

func commitHydrateCommitUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.CommitUrl, nil
}

func commitHydrateDeletions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Deletions, nil
}

func commitHydrateSignature(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Signature, nil
}

func commitHydrateTarballUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.TarballUrl, nil
}

func commitHydrateTreeUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.TreeUrl, nil
}

func commitHydrateCanSubscribe(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.CanSubscribe, nil
}

func commitHydrateSubscription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Subscription, nil
}

func commitHydrateZipballUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.ZipballUrl, nil
}

func commitHydrateMessageHeadline(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.MessageHeadline, nil
}

func commitHydrateStatus(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.Status, nil
}

func commitHydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	commit, err := extractCommitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return commit.NodeId, nil
}
