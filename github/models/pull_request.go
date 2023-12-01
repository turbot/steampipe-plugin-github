package models

import "github.com/shurcooL/githubv4"

type BasicPullRequest struct {
	Id                  int                                `graphql:"id: databaseId @include(if:$includePRId)" json:"id"`
	NodeId              string                             `graphql:"nodeId: id @include(if:$includePRNodeId)" json:"node_id"`
	Number              int                                `json:"number"`
	ActiveLockReason    githubv4.LockReason                `graphql:"activeLockReason @include(if:$includePRActiveLockReason)" json:"active_lock_reason"`
	Additions           int                                `graphql:"additions @include(if:$includePRAdditions)" json:"additions"`
	Author              Actor                              `graphql:"author @include(if:$includePRAuthor)" json:"author"`
	AuthorAssociation   githubv4.CommentAuthorAssociation  `graphql:"authorAssociation @include(if:$includePRAuthorAssociation)" json:"author_association"`
	BaseRefName         string                             `graphql:"baseRefName @include(if:$includePRBaseRefName)" json:"base_ref_name"`
	Body                string                             `graphql:"body @include(if:$includePRBody)" json:"body"`
	ChangedFiles        int                                `graphql:"changedFiles @include(if:$includePRChangedFiles)" json:"changed_files"`
	ChecksUrl           string                             `graphql:"checksUrl @include(if:$includePRChecksUrl)" json:"checks_url"`
	Closed              bool                               `graphql:"closed @include(if:$includePRClosed)" json:"closed"`
	ClosedAt            NullableTime                       `graphql:"closedAt @include(if:$includePRClosedAt)" json:"closed_at"`
	CreatedAt           NullableTime                       `graphql:"createdAt @include(if:$includePRCreatedAt)" json:"created_at"`
	CreatedViaEmail     bool                               `graphql:"createdViaEmail @include(if:$includePRCreatedViaEmail)" json:"created_via_email"`
	Deletions           int                                `graphql:"deletions @include(if:$includePRDeletions)" json:"deletions"`
	Editor              Actor                              `graphql:"editor @include(if:$includePREditor)" json:"editor"`
	HeadRefName         string                             `graphql:"headRefName @include(if:$includePRHeadRefName)" json:"head_ref_name"`
	HeadRefOid          string                             `graphql:"headRefOid @include(if:$includePRHeadRefOid)" json:"head_ref_oid"`
	IncludesCreatedEdit bool                               `graphql:"includesCreatedEdit @include(if:$includePRIncludesCreatedEdit)" json:"includes_created_edit"`
	IsCrossRepository   bool                               `graphql:"isCrossRepository @include(if:$includePRIsCrossRepository)" json:"is_cross_repository"`
	IsDraft             bool                               `graphql:"isDraft @include(if:$includePRIsDraft)" json:"is_draft"`
	IsReadByUser        bool                               `graphql:"isReadByUser: isReadByViewer @include(if:$includePRIsReadByUser)" json:"is_read_by_user"`
	LastEditedAt        NullableTime                       `graphql:"lastEditedAt @include(if:$includePRLastEditedAt)" json:"last_edited_at"`
	Locked              bool                               `graphql:"locked @include(if:$includePRLocked)" json:"locked"`
	MaintainerCanModify bool                               `graphql:"maintainerCanModify @include(if:$includePRMaintainerCanModify)" json:"maintainer_can_modify"`
	Mergeable           githubv4.MergeableState            `graphql:"mergeable @include(if:$includePRMergeable)" json:"mergeable"`
	Merged              bool                               `graphql:"merged @include(if:$includePRMerged)" json:"merged"`
	MergedAt            NullableTime                       `graphql:"mergedAt @include(if:$includePRMergedAt)" json:"merged_at"`
	MergedBy            Actor                              `graphql:"mergedBy @include(if:$includePRMergedBy)" json:"merged_by"`
	Milestone           Milestone                          `graphql:"milestone @include(if:$includePRMilestone)" json:"milestone"`
	Permalink           string                             `graphql:"permalink @include(if:$includePRPermalink)" json:"permalink"`
	PublishedAt         NullableTime                       `graphql:"publishedAt @include(if:$includePRPublishedAt)" json:"published_at"`
	RevertUrl           string                             `graphql:"revertUrl @include(if:$includePRRevertUrl)" json:"revert_url"`
	ReviewDecision      githubv4.PullRequestReviewDecision `graphql:"reviewDecision @include(if:$includePRReviewDecision)" json:"review_decision"`
	State               githubv4.PullRequestState          `graphql:"state @include(if:$includePRState)" json:"state"`
	Title               string                             `graphql:"title @include(if:$includePRTitle)" json:"title"`
	TotalCommentsCount  int                                `graphql:"totalCommentsCount @include(if:$includePRTotalCommentsCount)" json:"total_comments_count"`
	UpdatedAt           NullableTime                       `graphql:"updatedAt @include(if:$includePRUpdatedAt)" json:"updated_at"`
	Url                 string                             `graphql:"url @include(if:$includePRUrl)" json:"url"`
	Repo                struct {
		NameWithOwner string `json:"name_with_owner"`
	} `graphql:"repo: repository" json:"repo"`
	// CanBeRebased        bool                               `json:"can_be_rebased"`
}

type PullRequest struct {
	BasicPullRequest
	BaseRef             *BasicRef                            `graphql:"baseRef @include(if:$includePRBaseRef)" json:"base_ref,omitempty"`
	HeadRef             *BasicRef                            `graphql:"headRef @include(if:$includePRHeadRef)" json:"head_ref,omitempty"`
	MergeCommit         *BasicCommit                         `graphql:"mergeCommit @include(if:$includePRMergeCommit)" json:"merge_commit,omitempty"`
	SuggestedReviewers  []SuggestedReviewer                  `graphql:"suggestedReviewers @include(if:$includePRSuggested)" json:"suggested_reviewers"`
	CanApplySuggestion  bool                                 `graphql:"canApplySuggestion:viewerCanApplySuggestion @include(if:$includePRViewer)" json:"can_apply_suggestion"`
	CanClose            bool                                 `graphql:"canClose:viewerCanClose @include(if:$includePRViewer)" json:"can_close"`
	CanDeleteHeadRef    bool                                 `graphql:"canDeleteHeadRef:viewerCanDeleteHeadRef @include(if:$includePRViewer)" json:"can_delete_head_ref"`
	CanDisableAutoMerge bool                                 `graphql:"canDisableAutoMerge:viewerCanDisableAutoMerge @include(if:$includePRViewer)" json:"can_disable_auto_merge"`
	CanEditFiles        bool                                 `graphql:"canEditFiles:viewerCanEditFiles @include(if:$includePRViewer)" json:"can_edit_files"`
	CanEnableAutoMerge  bool                                 `graphql:"canEnableAutoMerge:viewerCanEnableAutoMerge @include(if:$includePRViewer)" json:"can_enable_auto_merge"`
	CanMergeAsAdmin     bool                                 `graphql:"canMergeAsAdmin:viewerCanMergeAsAdmin @include(if:$includePRViewer)" json:"can_merge_as_admin"`
	CanReact            bool                                 `graphql:"canReact:viewerCanReact @include(if:$includePRViewer)" json:"can_react"`
	CanReopen           bool                                 `graphql:"canReopen:viewerCanReopen @include(if:$includePRViewer)" json:"can_reopen"`
	CanSubscribe        bool                                 `graphql:"canSubscribe:viewerCanSubscribe @include(if:$includePRViewer)" json:"can_subscribe"`
	CanUpdate           bool                                 `graphql:"canUpdate:viewerCanUpdate @include(if:$includePRViewer)" json:"can_update"`
	CanUpdateBranch     bool                                 `graphql:"canUpdateBranch:viewerCanUpdateBranch @include(if:$includePRViewer)" json:"can_update_branch"`
	DidAuthor           bool                                 `graphql:"didAuthor:viewerDidAuthor @include(if:$includePRViewer)" json:"did_author"`
	CannotUpdateReasons []githubv4.CommentCannotUpdateReason `graphql:"cannotUpdateReasons: viewerCannotUpdateReasons @include(if:$includePRViewer)" json:"cannot_update_reasons"`
	Subscription        githubv4.SubscriptionState           `graphql:"subscription: viewerSubscription @include(if:$includePRViewer)" json:"subscription"`

	// Counts
	Commits        Count `graphql:"commits @include(if:$includePRCommitCount)" json:"commits"`
	ReviewRequests Count `graphql:"reviewRequests @include(if:$includePRReviewRequestCount)" json:"review_requests"`
	Reviews        Count `graphql:"reviews @include(if:$includePRReviewCount)" json:"reviews"`
	Labels         struct {
		TotalCount int
		Nodes      []Label
	} `graphql:"labels(first: 100) @include(if:$includePRLabels)" json:"labels"`
	Assignees struct {
		TotalCount int
		Nodes      []User
	} `graphql:"assignees(first: 10) @include(if:$includePRAssignees)" json:"assignees"`

	// ClosingIssueReferences [pageable]
	// Commits [pageable]
	// Files [pageable]
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

type PullRequestReview struct {
	Id                        int                               `graphql:"id: databaseId @include(if:$includePRReviewId)" json:"id"`
	NodeId                    string                            `graphql:"nodeId: id @include(if:$includePRReviewNodeId)" json:"node_id"`
	Author                    Actor                             `graphql:"author @include(if:$includePRReviewAuthor)" json:"author"`
	AuthorAssociation         githubv4.CommentAuthorAssociation `graphql:"authorAssociation @include(if:$includePRReviewAuthorAssociation)" json:"author_association"`
	AuthorCanPushToRepository bool                              `graphql:"authorCanPushToRepository @include(if:$includePRReviewAuthorCanPushToRepository)" json:"author_can_push_to_repository"`
	State                     string                            `graphql:"state @include(if:$includePRReviewState)" json:"state"`
	Body                      string                            `graphql:"body @include(if:$includePRReviewBody)" json:"body"`
	Url                       string                            `graphql:"url @include(if:$includePRReviewUrl)" json:"html_url"`
	SubmittedAt               NullableTime                      `graphql:"submittedAt @include(if:$includePRReviewSubmittedAt)" json:"submitted_at"`
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
