package models

type RepositoryInteractionAbility struct {
	ExpiresAt NullableTime `json:"expires_at,omitempty"`
	Limit     string       `json:"repository_interaction_limit,omitempty"`
	Origin    string       `json:"repository_interaction_limit_origin,omitempty"`
}
