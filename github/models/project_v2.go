package models

// ProjectV2Owner represents the owner of a project, which can be an Organization or User.
type ProjectV2Owner struct {
	TypeName     string `graphql:"type: __typename" json:"type"`
	Organization struct {
		Id    int    `graphql:"id: databaseId" json:"id"`
		Login string `json:"login"`
	} `graphql:"... on Organization" json:"organization,omitempty"`
	User struct {
		Id    int    `graphql:"id: databaseId" json:"id"`
		Login string `json:"login"`
	} `graphql:"... on User" json:"user,omitempty"`
}

// ProjectV2StatusUpdate represents a single status update on a project.
type ProjectV2StatusUpdate struct {
	Id         string       `graphql:"id: fullDatabaseId" json:"id"`
	NodeId     string       `graphql:"nodeId: id" json:"node_id"`
	Status     string       `json:"status"`
	Body       string       `json:"body"`
	StartDate  string       `json:"start_date"`
	TargetDate string       `json:"target_date"`
	CreatedAt  NullableTime `json:"created_at"`
	UpdatedAt  NullableTime `json:"updated_at"`
	Creator    Actor        `json:"creator"`
}

type ProjectV2 struct {
	Id                  string         `graphql:"id: fullDatabaseId @include(if:$includeId)" json:"id"`
	NodeId              string         `graphql:"nodeId: id @include(if:$includeNodeId)" json:"node_id"`
	Number              int            `json:"number"`
	Owner               ProjectV2Owner `graphql:"owner @include(if:$includeOwner)" json:"owner,omitempty"`
	Creator             Actor          `graphql:"creator @include(if:$includeCreator)" json:"creator,omitempty"`
	Title               string         `graphql:"title @include(if:$includeTitle)" json:"title"`
	Description         string         `graphql:"description: shortDescription @include(if:$includeDescription)" json:"description"`
	IsPublic            bool           `graphql:"public @include(if:$includeIsPublic)" json:"public"`
	ClosedAt            NullableTime   `graphql:"closedAt @include(if:$includeClosedAt)" json:"closed_at"`
	CreatedAt           NullableTime   `graphql:"createdAt @include(if:$includeCreatedAt)" json:"created_at"`
	UpdatedAt           NullableTime   `graphql:"updatedAt @include(if:$includeUpdatedAt)" json:"updated_at"`
	Closed              bool           `graphql:"closed @include(if:$includeState)" json:"closed"`
	LatestStatusUpdate struct {
		Nodes []ProjectV2StatusUpdate
	} `graphql:"statusUpdates(last: 1) @include(if:$includeLatestStatusUpdate)" json:"latest_status_update"`
	IsTemplate bool `graphql:"template @include(if:$includeIsTemplate)" json:"template"`
}
