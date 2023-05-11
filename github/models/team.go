package models

import (
	"time"
)

type Team struct {
	basicIdentifiers
	AvatarUrl      string
	CombinedSlug   string
	CreatedAt      time.Time
	Description    string
	DiscussionsUrl string
	EditTeamUrl    string
	MembersUrl     string
	NewTeamUrl     string
	Organization   BasicOrganization
	ParentTeam     struct {
		basicIdentifiers
		Slug string `json:"slug,omitempty"`
	}
	Privacy         string
	RepositoriesUrl string
	Slug            string
	TeamsUrl        string
	UpdatedAt       time.Time
	Url             string
	CanAdminister   bool   `graphql:"canAdminister: viewerCanAdminister"`
	CanSubscribe    bool   `graphql:"canSubscribe: viewerCanSubscribe"`
	Subscription    string `graphql:"subscription: viewerSubscription"`
}

type TeamWithCounts struct {
	Team
	Ancestors struct {
		TotalCount int
	}
	ChildTeams struct {
		TotalCount int
	}
	Discussions struct {
		TotalCount int
	}
	Invitations struct {
		TotalCount int
	}
	Members struct {
		TotalCount int
	}
	ProjectsV2 struct {
		TotalCount int
	}
	Repositories struct {
		TotalCount int
	}
}
