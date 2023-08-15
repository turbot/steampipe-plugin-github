package models

type Actor struct {
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Url       string `json:"url"`
}

type GitActor struct {
	AvatarUrl string       `json:"avatar_url"`
	Date      NullableTime `json:"date"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	User      BasicUser    `json:"user"`
}
