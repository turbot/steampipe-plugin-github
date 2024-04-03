package models

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type Team struct {
	basicIdentifiers
	AvatarUrl      string            `graphql:"avatarUrl @include(if:$includeTeamAvatarUrl)" json:"avatar_url"`
	CombinedSlug   string            `graphql:"combinedSlug @include(if:$includeTeamCombinedSlug)" json:"combined_slug"`
	CreatedAt      time.Time         `graphql:"createdAt @include(if:$includeTeamCreatedAt)" json:"created_at"`
	Description    string            `graphql:"description @include(if:$includeTeamDescription)" json:"description"`
	DiscussionsUrl string            `graphql:"discussionsUrl @include(if:$includeTeamDiscussionsUrl)" json:"discussions_url"`
	EditTeamUrl    string            `graphql:"editTeamUrl @include(if:$includeTeamEditTeamUrl)" json:"edit_team_url"`
	MembersUrl     string            `graphql:"membersUrl @include(if:$includeTeamMembersUrl)" json:"members_url"`
	NewTeamUrl     string            `graphql:"newTeamUrl @include(if:$includeTeamNewTeamUrl)" json:"new_team_url"`
	Organization   BasicOrganization `json:"organization"`
	ParentTeam     struct {
		basicIdentifiers
		Slug string `json:"slug,omitempty"`
	} `graphql:"parentTeam @include(if:$includeTeamParentTeam)" json:"parent_team"`
	Privacy         string    `graphql:"privacy @include(if:$includeTeamPrivacy)" json:"privacy"`
	RepositoriesUrl string    `graphql:"repositoriesUrl @include(if:$includeTeamRepositoriesUrl)" json:"repositories_url"`
	Slug            string    `json:"slug"`
	TeamsUrl        string    `graphql:"teamsUrl @include(if:$includeTeamTeamsUrl)" json:"teams_url"`
	UpdatedAt       time.Time `graphql:"updatedAt @include(if:$includeTeamUpdatedAt)" json:"updated_at"`
	Url             string    `graphql:"url @include(if:$includeTeamUrl)" json:"url"`
	CanAdminister   bool      `graphql:"canAdminister: viewerCanAdminister @include(if:$includeTeamCanAdminister)" json:"can_administer"`
	CanSubscribe    bool      `graphql:"canSubscribe: viewerCanSubscribe @include(if:$includeTeamCanSubscribe)" json:"can_subscribe"`
	Subscription    string    `graphql:"subscription: viewerSubscription @include(if:$includeTeamSubscription)" json:"subscription"`
}

type TeamWithCounts struct {
	Team
	Ancestors    Count `graphql:"ancestors @include(if:$includeTeamAncestors)" json:"ancestors"`
	ChildTeams   Count `graphql:"childTeams @include(if:$includeTeamChildTeams)" json:"child_teams"`
	Discussions  Count `graphql:"discussions @include(if:$includeTeamDiscussions)" json:"discussions"`
	Invitations  Count `graphql:"invitations @include(if:$includeTeamInvitations)" json:"invitations"`
	Members      Count `graphql:"members @include(if:$includeTeamMembers)" json:"members"`
	ProjectsV2   Count `graphql:"projectsV2 @include(if:$includeTeamProjectsV2)" json:"projects_v2"`
	Repositories Count `graphql:"repositories @include(if:$includeTeamRepositories)" json:"repositories"`
}

type TeamMemberWithRole struct {
	Role githubv4.TeamMemberRole `json:"role"`
	Node User
}

type TeamRepositoryWithPermission struct {
	Permission githubv4.RepositoryPermission `json:"permission"`
	Node       Repository
}
