package models

type Discussion struct {
	Id        int          `graphql:"id: databaseId" json:"id"`
	NodeId    string       `graphql:"nodeId: id" json:"node_id"`
	Number    int          `graphql:"number" json:"number"`
	Title     string       `graphql:"title" json:"title"`
	Url       string       `graphql:"url" json:"url"`
	CreatedAt NullableTime `graphql:"createdAt" json:"created_at"`
	UpdatedAt NullableTime `graphql:"updatedAt" json:"updated_at"`
	Author    Actor        `graphql:"author" json:"author"`
	Category  struct {
		Name string `graphql:"name" json:"name"`
	} `graphql:"category" json:"category"`
	Answer   *DiscussionComment `graphql:"answer" json:"answer"`
	Comments struct {
		TotalCount int
		Nodes      []DiscussionComment
	} `graphql:"comments(first: 10)" json:"comments"`
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
