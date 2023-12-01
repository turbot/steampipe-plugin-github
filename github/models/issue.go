package models

import "github.com/shurcooL/githubv4"

type Issue struct {
	Id                      int                                  `graphql:"id: databaseId @include(if:$includeIssueId)" json:"id"`
	NodeId                  string                               `graphql:"nodeId: id @include(if:$includeIssueNodeId)" json:"node_id"`
	Number                  int                                  `json:"number"`
	ActiveLockReason        githubv4.LockReason                  `graphql:"activeLockReason @include(if:$includeIssueActiveLockReason)" json:"active_lock_reason"`
	Author                  Actor                                `graphql:"author @include(if:$includeIssueAuthor)" json:"author"`
	AuthorAssociation       githubv4.CommentAuthorAssociation    `graphql:"authorAssociation @include(if:$includeIssueAuthorAssociation)" json:"author_association"`
	Body                    string                               `graphql:"body @include(if:$includeIssueBody)" json:"body"`
	BodyUrl                 string                               `graphql:"bodyUrl @include(if:$includeIssueBodyUrl)" json:"body_url"`
	Closed                  bool                                 `graphql:"closed @include(if:$includeIssueClosed)" json:"closed"`
	ClosedAt                NullableTime                         `graphql:"closedAt @include(if:$includeIssueClosedAt)" json:"closed_at"`
	CreatedAt               NullableTime                         `graphql:"createdAt @include(if:$includeIssueCreatedAt)" json:"created_at"`
	CreatedViaEmail         bool                                 `graphql:"createdViaEmail @include(if:$includeIssueCreatedViaEmail)" json:"created_via_email"`
	Editor                  Actor                                `graphql:"editor @include(if:$includeIssueEditor)" json:"editor"`
	FullDatabaseId          string                               `graphql:"fullDatabaseId @include(if:$includeIssueFullDatabaseId)" json:"full_database_id"`
	IncludesCreatedEdit     bool                                 `graphql:"includesCreatedEdit @include(if:$includeIssueIncludesCreatedEdit)" json:"includes_created_edit"`
	IsPinned                bool                                 `graphql:"isPinned @include(if:$includeIssueIsPinned)" json:"is_pinned"`
	IsReadByUser            bool                                 `graphql:"isReadByUser: isReadByViewer @include(if:$includeIssueIsReadByUser)" json:"is_read_by_user"`
	LastEditedAt            NullableTime                         `graphql:"lastEditedAt @include(if:$includeIssueLastEditedAt)" json:"last_edited_at"`
	Locked                  bool                                 `graphql:"locked @include(if:$includeIssueLocked)" json:"locked"`
	Milestone               Milestone                            `graphql:"milestone @include(if:$includeIssueMilestone)" json:"milestone"`
	PublishedAt             NullableTime                         `graphql:"publishedAt @include(if:$includeIssuePublishedAt)" json:"published_at"`
	State                   githubv4.IssueState                  `graphql:"state @include(if:$includeIssueState)" json:"state"`
	StateReason             githubv4.IssueStateReason            `graphql:"stateReason @include(if:$includeIssueStateReason)" json:"state_reason"`
	Title                   string                               `graphql:"title @include(if:$includeIssueTitle)" json:"title"`
	UpdatedAt               NullableTime                         `graphql:"updatedAt @include(if:$includeIssueUpdatedAt)" json:"updated_at"`
	Url                     string                               `graphql:"url @include(if:$includeIssueUrl)" json:"url"`
	UserCanClose            bool                                 `graphql:"userCanClose: viewerCanClose @include(if:$includeIssueViewer)" json:"user_can_close"`
	UserCanReact            bool                                 `graphql:"userCanReact: viewerCanReact @include(if:$includeIssueViewer)" json:"user_can_react"`
	UserCanReopen           bool                                 `graphql:"userCanReopen: viewerCanReopen @include(if:$includeIssueViewer)" json:"user_can_reopen"`
	UserCanSubscribe        bool                                 `graphql:"userCanSubscribe: viewerCanSubscribe @include(if:$includeIssueViewer)" json:"user_can_subscribe"`
	UserCanUpdate           bool                                 `graphql:"userCanUpdate: viewerCanUpdate @include(if:$includeIssueViewer)" json:"user_can_update"`
	UserCannotUpdateReasons []githubv4.CommentCannotUpdateReason `graphql:"userCannotUpdateReasons: viewerCannotUpdateReasons @include(if:$includeIssueViewer)" json:"user_cannot_update_reasons"`
	UserDidAuthor           bool                                 `graphql:"userDidAuthor: viewerDidAuthor @include(if:$includeIssueViewer)" json:"user_did_author"`
	UserSubscription        githubv4.SubscriptionState           `graphql:"userSubscription: viewerSubscription @include(if:$includeIssueViewer)" json:"user_subscription"`
	Comments                Count                                `graphql:"comments @include(if:$includeIssueCommentCount)" json:"comments"`
	Labels                  struct {
		TotalCount int
		Nodes      []Label
	} `graphql:"labels(first: 100) @include(if:$includeIssueLabels)" json:"labels"`
	Repo struct {
		NameWithOwner string `json:"name_with_owner"`
	} `graphql:"repo: repository" json:"repo"`
	Assignees struct {
		TotalCount int
		Nodes      []User
	} `graphql:"assignees(first: 10) @include(if:$includeIssueAssignees)" json:"assignees"`
}

type IssueTemplate struct {
	About    string `json:"about"`
	Body     string `json:"body"`
	Filename string `json:"filename"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	// Assignees [pageable]
	// Labels [pageable]
}

type IssueComment struct {
	Id                  int                                  `graphql:"id: databaseId @include(if:$includeIssueCommentId)" json:"id"`
	NodeId              string                               `graphql:"nodeId: id @include(if:$includeIssueCommentNodeId)" json:"node_id"`
	Author              Actor                                `graphql:"author @include(if:$includeIssueCommentAuthor)" json:"author"`
	AuthorAssociation   githubv4.CommentAuthorAssociation    `graphql:"authorAssociation @include(if:$includeIssueCommentAuthorAssociation)" json:"author_association"`
	Body                string                               `graphql:"body @include(if:$includeIssueCommentBody)" json:"body"`
	BodyText            string                               `graphql:"bodyText @include(if:$includeIssueCommentBodyText)" json:"body_text"`
	CreatedAt           NullableTime                         `graphql:"createdAt @include(if:$includeIssueCommentCreatedAt)" json:"created_at"`
	CreatedViaEmail     bool                                 `graphql:"createdViaEmail @include(if:$includeIssueCommentCreatedViaEmail)" json:"created_via_email"`
	Editor              Actor                                `graphql:"editor @include(if:$includeIssueCommentEditor)" json:"editor"`
	IncludesCreatedEdit bool                                 `graphql:"includesCreatedEdit @include(if:$includeIssueCommentIncludesCreatedEdit)" json:"includes_created_edit"`
	IsMinimized         bool                                 `graphql:"isMinimized @include(if:$includeIssueCommentIsMinimized)" json:"is_minimized"`
	LastEditedAt        NullableTime                         `graphql:"lastEditedAt @include(if:$includeIssueCommentLastEditedAt)" json:"last_edited_at"`
	MinimizedReason     string                               `graphql:"minimizedReason @include(if:$includeIssueCommentMinimizedReason)" json:"minimized_reason"`
	PublishedAt         NullableTime                         `graphql:"publishedAt @include(if:$includeIssueCommentPublishedAt)" json:"published_at"`
	UpdatedAt           NullableTime                         `graphql:"updatedAt @include(if:$includeIssueCommentUpdatedAt)" json:"updated_at"`
	Url                 string                               `graphql:"url @include(if:$includeIssueCommentUrl)" json:"url"`
	CanDelete           bool                                 `graphql:"canDelete: viewerCanDelete @include(if:$includeIssueCommentViewer)" json:"can_delete"`
	CanMinimize         bool                                 `graphql:"canMinimize: viewerCanMinimize @include(if:$includeIssueCommentViewer)" json:"can_minimize"`
	CanReact            bool                                 `graphql:"canReact: viewerCanReact @include(if:$includeIssueCommentViewer)" json:"can_react"`
	CanUpdate           bool                                 `graphql:"canUpdate: viewerCanUpdate @include(if:$includeIssueCommentViewer)" json:"can_update"`
	CannotUpdateReasons []githubv4.CommentCannotUpdateReason `graphql:"cannotUpdateReasons: viewerCannotUpdateReasons @include(if:$includeIssueCommentViewer)" json:"cannot_update_reasons"`
	DidAuthor           bool                                 `graphql:"didAuthor: viewerDidAuthor @include(if:$includeIssueCommentViewer)" json:"did_author"`
}
