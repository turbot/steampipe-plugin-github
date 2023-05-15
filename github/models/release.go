package models

type Release struct {
	basicIdentifiers
	Author       User         `json:"author"`
	CreatedAt    NullableTime `json:"created_at"`
	Description  string       `json:"description"`
	IsDraft      bool         `json:"is_draft"`
	IsLatest     bool         `json:"is_latest"`
	IsPrerelease bool         `json:"is_prerelease"`
	PublishedAt  NullableTime `json:"published_at"`
	Tag          BasicRef     `json:"tag"`
	TagCommit    Commit       `json:"tag_commit"`
	TagName      string       `json:"tag_name"`
	UpdatedAt    NullableTime `json:"updated_at"`
	Url          string       `json:"url"`
	CanReact     bool         `graphql:"canReact: viewerCanReact" json:"can_react"`
	// Mentions [pageable]
	// Reactions [pageable]
	// ReleaseAssets [pageable]
}
