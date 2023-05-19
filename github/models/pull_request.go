package models

import "github.com/shurcooL/githubv4"

type BasicPullRequest struct {
	Id                  int                                `graphql:"id: databaseId" json:"id"`
	NodeId              string                             `graphql:"nodeId: id" json:"node_id"`
	Number              int                                `json:"number"`
	ActiveLockReason    githubv4.LockReason                `json:"active_lock_reason"`
	Additions           int                                `json:"additions"`
	Author              Actor                              `json:"author"`
	AuthorAssociation   githubv4.CommentAuthorAssociation  `json:"author_association"`
	BaseRefName         string                             `json:"base_ref_name"`
	Body                string                             `json:"body"`
	ChangedFiles        int                                `json:"changed_files"`
	ChecksUrl           string                             `json:"checks_url"`
	Closed              bool                               `json:"closed"`
	ClosedAt            NullableTime                       `json:"closed_at"`
	CreatedAt           NullableTime                       `json:"created_at"`
	CreatedViaEmail     bool                               `json:"created_via_email"`
	Deletions           int                                `json:"deletions"`
	Editor              Actor                              `json:"editor"`
	HeadRefName         string                             `json:"head_ref_name"`
	HeadRefOid          string                             `json:"head_ref_oid"`
	IncludesCreatedEdit bool                               `json:"includes_created_edit"`
	IsCrossRepository   bool                               `json:"is_cross_repository"`
	IsDraft             bool                               `json:"is_draft"`
	IsReadByUser        bool                               `graphql:"isReadByUser: isReadByViewer" json:"is_read_by_user"`
	LastEditedAt        NullableTime                       `json:"last_edited_at"`
	Locked              bool                               `json:"locked"`
	MaintainerCanModify bool                               `json:"maintainer_can_modify"`
	Mergeable           githubv4.MergeableState            `json:"mergeable"`
	Merged              bool                               `json:"merged"`
	MergedAt            NullableTime                       `json:"merged_at"`
	MergedBy            Actor                              `json:"merged_by"`
	Milestone           Milestone                          `json:"milestone"`
	Permalink           string                             `json:"permalink"`
	PublishedAt         NullableTime                       `json:"published_at"`
	RevertUrl           string                             `json:"revert_url"`
	ReviewDecision      githubv4.PullRequestReviewDecision `json:"review_decision"`
	State               githubv4.PullRequestState          `json:"state"`
	Title               string                             `json:"title"`
	TotalCommentsCount  int                                `json:"total_comments_count"`
	UpdatedAt           NullableTime                       `json:"updated_at"`
	Url                 string                             `json:"url"`
	Repo                struct {
		NameWithOwner string `json:"name_with_owner"`
	} `graphql:"repo: repository" json:"repo"`
}

type PullRequest struct {
	BasicPullRequest
	BaseRef             BasicRef                             `json:"base_ref"`
	HeadRef             BasicRef                             `json:"head_ref"`
	MergeCommit         Commit                               `json:"merge_commit"`
	SuggestedReviewers  []SuggestedReviewer                  `json:"suggested_reviewers"`
	CanApplySuggestion  bool                                 `graphql:"canApplySuggestion:viewerCanApplySuggestion" json:"can_apply_suggestion"`
	CanClose            bool                                 `graphql:"canClose:viewerCanClose" json:"can_close"`
	CanDeleteHeadRef    bool                                 `graphql:"canDeleteHeadRef:viewerCanDeleteHeadRef" json:"can_delete_head_ref"`
	CanDisableAutoMerge bool                                 `graphql:"canDisableAutoMerge:viewerCanDisableAutoMerge" json:"can_disable_auto_merge"`
	CanEditFiles        bool                                 `graphql:"canEditFiles:viewerCanEditFiles" json:"can_edit_files"`
	CanEnableAutoMerge  bool                                 `graphql:"canEnableAutoMerge:viewerCanEnableAutoMerge" json:"can_enable_auto_merge"`
	CanMergeAsAdmin     bool                                 `graphql:"canMergeAsAdmin:viewerCanMergeAsAdmin" json:"can_merge_as_admin"`
	CanReact            bool                                 `graphql:"canReact:viewerCanReact" json:"can_react"`
	CanReopen           bool                                 `graphql:"canReopen:viewerCanReopen" json:"can_reopen"`
	CanSubscribe        bool                                 `graphql:"canSubscribe:viewerCanSubscribe" json:"can_subscribe"`
	CanUpdate           bool                                 `graphql:"canUpdate:viewerCanUpdate" json:"can_update"`
	CanUpdateBranch     bool                                 `graphql:"canUpdateBranch:viewerCanUpdateBranch" json:"can_update_branch"`
	DidAuthor           bool                                 `graphql:"didAuthor:viewerDidAuthor" json:"did_author"`
	CannotUpdateReasons []githubv4.CommentCannotUpdateReason `graphql:"cannotUpdateReasons: viewerCannotUpdateReasons" json:"cannot_update_reasons"`
	Subscription        githubv4.SubscriptionState           `graphql:"subscription: viewerSubscription" json:"subscription"`

	// Assignees [pageable]
	// ClosingIssueReferences [pageable]
	// Comments [pageable]
	// Commits [pageable]
	// Files [pageable]
	// Labels [pageable]
	// LatestOpinionatedReviews [pageable]
	// LatestReviews [pageable]
	// Participants [pageable]
	// ProjectCards [pageable]
	// ProjectItems [pageable]
	// ProjectV2 [find by number[
	// ProjectsV2 [pageable]
	// Reactions [pageable]
	// ReviewRequests [pageable]
	// ReviewThreads [pageable]
	// Reviews [pageable]
	// TimelineItems [pageable]
}

type SuggestedReviewer struct {
	IsAuthor    bool      `json:"is_author"`
	IsCommenter bool      `json:"is_commenter"`
	Reviewer    BasicUser `json:"reviewer"`
}

type PullRequestTemplate struct {
	Filename string `json:"filename"`
	Body     string `json:"body"`
}
