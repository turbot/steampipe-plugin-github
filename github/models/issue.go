package models

import "github.com/shurcooL/githubv4"

type Issue struct {
	Id                      int                                  `graphql:"id: databaseId" json:"id"`
	NodeId                  string                               `graphql:"nodeId: id" json:"node_id"`
	Number                  int                                  `json:"number"`
	ActiveLockReason        githubv4.LockReason                  `json:"active_lock_reason"`
	Author                  Actor                                `json:"author"`
	AuthorAssociation       githubv4.CommentAuthorAssociation    `json:"author_association"`
	Body                    string                               `json:"body"`
	BodyUrl                 string                               `json:"body_url"`
	Closed                  bool                                 `json:"closed"`
	ClosedAt                NullableTime                         `json:"closed_at"`
	CreatedAt               NullableTime                         `json:"created_at"`
	CreatedViaEmail         bool                                 `json:"created_via_email"`
	Editor                  Actor                                `json:"editor"`
	FullDatabaseId          string                               `json:"full_database_id"`
	IncludesCreatedEdit     bool                                 `json:"includes_created_edit"`
	IsPinned                bool                                 `json:"is_pinned"`
	IsReadByUser            bool                                 `graphql:"isReadByUser: isReadByViewer" json:"is_read_by_user"`
	LastEditedAt            NullableTime                         `json:"last_edited_at"`
	Locked                  bool                                 `json:"locked"`
	Milestone               Milestone                            `json:"milestone"`
	PublishedAt             NullableTime                         `json:"published_at"`
	State                   githubv4.IssueState                  `json:"state"`
	StateReason             githubv4.IssueStateReason            `json:"state_reason"`
	Title                   string                               `json:"title"`
	UpdatedAt               NullableTime                         `json:"updated_at"`
	Url                     string                               `json:"url"`
	UserCanClose            bool                                 `graphql:"userCanClose: viewerCanClose" json:"user_can_close"`
	UserCanReact            bool                                 `graphql:"userCanReact: viewerCanReact" json:"user_can_react"`
	UserCanReopen           bool                                 `graphql:"userCanReopen: viewerCanReopen" json:"user_can_reopen"`
	UserCanSubscribe        bool                                 `graphql:"userCanSubscribe: viewerCanSubscribe" json:"user_can_subscribe"`
	UserCanUpdate           bool                                 `graphql:"userCanUpdate: viewerCanUpdate" json:"user_can_update"`
	UserCannotUpdateReasons []githubv4.CommentCannotUpdateReason `graphql:"userCannotUpdateReasons: viewerCannotUpdateReasons" json:"user_cannot_update_reasons"`
	UserDidAuthor           bool                                 `graphql:"userDidAuthor: viewerDidAuthor" json:"user_did_author"`
	UserSubscription        githubv4.SubscriptionState           `graphql:"userSubscription: viewerSubscription" json:"user_subscription"`
	Comments                Count                                `json:"comments"`
	Repo                    struct {
		NameWithOwner string `json:"name_with_owner"`
	} `graphql:"repo: repository" json:"repo"`

	// Assignees [pageable]
	// Comments [pageable]
	// Labels [pageable]
	// LinkedBranches [pageable]
	// Participants [pageable]
	// ProjectCards [pageable]
	// ProjectItems [pageable]
	// ProjectV2 [find by number]
	// ProjectsV2 [pageable]
	// Reactions [pageable]
	// TimelineItems [pageable]
	// TrackedInIssues [pageable]
	// TrackedIssues [pageable]
	// UserContentEdits [pageable]
}

// IssueWithRepository should not be nested under the repository, else circular reference will cause long wait and eventual error.
type IssueWithRepository struct {
	Issue
	Repository Repository `json:"repository"`
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
