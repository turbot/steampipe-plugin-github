package models

type Discussion struct {
	Id        int          `graphql:"id: databaseId @include(if:$includeDiscussionId)" json:"id"`
	NodeId    string       `graphql:"nodeId: id @include(if:$includeDiscussionNodeId)" json:"node_id"`
	Number    int          `graphql:"number" json:"number"`
	Title     string       `graphql:"title" json:"title"`
	Url       string       `graphql:"url" json:"url"`
	CreatedAt NullableTime `graphql:"createdAt @include(if:$includeDiscussionCreatedAt)" json:"created_at"`
	UpdatedAt NullableTime `graphql:"updatedAt @include(if:$includeDiscussionUpdatedAt)" json:"updated_at"`
	Author    Actor        `graphql:"author @include(if:$includeDiscussionAuthor)" json:"author"`
	Category  struct {
		Name string `graphql:"name" json:"name"`
	} `graphql:"category @include(if:$includeDiscussionCategory)" json:"category"`
	Answer   *DiscussionComment `graphql:"answer @include(if:$includeDiscussionAnswer)" json:"answer"`
	Comments struct {
		TotalCount int
		Nodes      []DiscussionComment
	} `graphql:"comments(first: 10) @include(if:$includeDiscussionComments)" json:"comments"`
}

type DiscussionComments struct {
	Comments struct {
		PageInfo   PageInfo
		TotalCount int
		Nodes      []DiscussionComment
	} `graphql:"comments(first: $pageSize, after: $cursor)"`
}

type DiscussionCommentReplies struct {
	Replies struct {
		PageInfo   PageInfo
		TotalCount int
		Nodes      []DiscussionComment
	} `graphql:"replies(first: $pageSize, after: $cursor)"`
}

type DiscussionComment struct {
	Id        int          `graphql:"id: databaseId" json:"id"`
	NodeId    string       `graphql:"nodeId: id" json:"node_id"`
	Author    Actor        `graphql:"author" json:"author"`
	Body      string       `graphql:"body" json:"body"`
	BodyText  string       `graphql:"bodyText" json:"body_text"`
	CreatedAt NullableTime `graphql:"createdAt" json:"created_at"`
	UpdatedAt NullableTime `graphql:"updatedAt" json:"updated_at"`
	IsAnswer  bool         `graphql:"isAnswer" json:"is_answer"`
}
