package models

import "github.com/shurcooL/githubv4"

// ProjectV2ItemArchivedState isn't defined in the pinned githubv4 version. The
// type name must match the GraphQL enum name, since githubv4 uses it when
// declaring query variables.
type ProjectV2ItemArchivedState string

const (
	ProjectV2ItemArchivedStateArchived    ProjectV2ItemArchivedState = "ARCHIVED"
	ProjectV2ItemArchivedStateNotArchived ProjectV2ItemArchivedState = "NOT_ARCHIVED"
)

type ProjectV2Item struct {
	Id          string                     `graphql:"id: fullDatabaseId" json:"id"`
	NodeId      string                     `graphql:"nodeId: id" json:"node_id"`
	Type        githubv4.ProjectV2ItemType `json:"type"`
	IsArchived  bool                       `json:"is_archived"`
	CreatedAt   NullableTime               `json:"created_at"`
	UpdatedAt   NullableTime               `json:"updated_at"`
	Creator     Actor                      `graphql:"creator @include(if:$includeItemCreator)" json:"creator"`
	Content     ProjectV2ItemContent       `graphql:"content @include(if:$includeItemContent)" json:"content"`
	FieldValues struct {
		Nodes []ProjectV2ItemFieldValue
	} `graphql:"fieldValues(first: 50) @include(if:$includeItemFieldValues)" json:"field_values"`
}

// ProjectV2ItemContent holds the item's content, one struct per possible type.
// The GraphQL client fills every fragment that shares a JSON key with the
// response, so read the fragment that matches the item's type.
type ProjectV2ItemContent struct {
	Issue struct {
		Id         int    `graphql:"id: databaseId"`
		NodeId     string `graphql:"nodeId: id"`
		Number     int
		Title      string
		Body       string `graphql:"body @include(if:$includeItemBody)"`
		Url        string
		State      githubv4.IssueState `graphql:"issueState: state"`
		Closed     bool
		ClosedAt   NullableTime
		CreatedAt  NullableTime
		UpdatedAt  NullableTime
		Author     Actor `graphql:"author @include(if:$includeItemAuthor)"`
		Repository struct {
			NameWithOwner string
		}
		Assignees ProjectV2ItemAssignees `graphql:"assignees(first: 10) @include(if:$includeItemAssignees)"`
		Labels    ProjectV2ItemLabels    `graphql:"labels(first: 100) @include(if:$includeItemLabels)"`
		Milestone Milestone              `graphql:"milestone @include(if:$includeItemMilestone)"`
	} `graphql:"... on Issue"`
	PullRequest struct {
		Id         int    `graphql:"id: databaseId"`
		NodeId     string `graphql:"nodeId: id"`
		Number     int
		Title      string
		Body       string `graphql:"body @include(if:$includeItemBody)"`
		Url        string
		State      githubv4.PullRequestState `graphql:"pullRequestState: state"`
		Closed     bool
		ClosedAt   NullableTime
		CreatedAt  NullableTime
		UpdatedAt  NullableTime
		Author     Actor `graphql:"author @include(if:$includeItemAuthor)"`
		Repository struct {
			NameWithOwner string
		}
		Assignees ProjectV2ItemAssignees `graphql:"assignees(first: 10) @include(if:$includeItemAssignees)"`
		Labels    ProjectV2ItemLabels    `graphql:"labels(first: 100) @include(if:$includeItemLabels)"`
		Milestone Milestone              `graphql:"milestone @include(if:$includeItemMilestone)"`
	} `graphql:"... on PullRequest"`
	DraftIssue struct {
		NodeId    string `graphql:"nodeId: id"`
		Title     string
		Body      string `graphql:"body @include(if:$includeItemBody)"`
		CreatedAt NullableTime
		UpdatedAt NullableTime
		Creator   Actor                  `graphql:"creator @include(if:$includeItemAuthor)"`
		Assignees ProjectV2ItemAssignees `graphql:"assignees(first: 10) @include(if:$includeItemAssignees)"`
	} `graphql:"... on DraftIssue"`
}

type ProjectV2ItemAssignees struct {
	TotalCount int
	Nodes      []BasicUser
}

type ProjectV2ItemLabels struct {
	Nodes []Label
}

// ProjectV2ItemFieldValue only selects the value types of custom fields. Built-in
// fields such as Assignees and Labels come back with just their type name.
type ProjectV2ItemFieldValue struct {
	TypeName string `graphql:"typeName: __typename"`
	Text     struct {
		Field ProjectV2FieldRef
		Text  string
	} `graphql:"... on ProjectV2ItemFieldTextValue"`
	Number struct {
		Field  ProjectV2FieldRef
		Number float64
	} `graphql:"... on ProjectV2ItemFieldNumberValue"`
	Date struct {
		Field ProjectV2FieldRef
		Date  string
	} `graphql:"... on ProjectV2ItemFieldDateValue"`
	SingleSelect struct {
		Field ProjectV2FieldRef
		Name  string
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	MultiSelect struct {
		Field   ProjectV2FieldRef
		Options []struct {
			Name string
		}
	} `graphql:"... on ProjectV2ItemFieldMultiSelectValue"`
	Iteration struct {
		Field       ProjectV2FieldRef
		Title       string
		StartDate   string
		Duration    int
		IterationId string
	} `graphql:"... on ProjectV2ItemFieldIterationValue"`
	IssueField struct {
		Field           ProjectV2FieldRef
		IssueFieldValue struct {
			Text struct {
				StringValue string `graphql:"stringValue: value"`
			} `graphql:"... on IssueFieldTextValue"`
			Date struct {
				StringValue string `graphql:"stringValue: value"`
			} `graphql:"... on IssueFieldDateValue"`
			SingleSelect struct {
				StringValue string `graphql:"stringValue: value"`
			} `graphql:"... on IssueFieldSingleSelectValue"`
			MultiSelect struct {
				StringValue string `graphql:"stringValue: value"`
			} `graphql:"... on IssueFieldMultiSelectValue"`
			Number struct {
				NumberValue *float64 `graphql:"numberValue: value"`
			} `graphql:"... on IssueFieldNumberValue"`
		}
	} `graphql:"... on ProjectV2ItemIssueFieldValue"`
}

type ProjectV2FieldRef struct {
	Common struct {
		Name     string
		DataType string
	} `graphql:"... on ProjectV2FieldCommon"`
}
