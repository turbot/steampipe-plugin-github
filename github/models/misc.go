package models

type NameSlug struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type NameLogin struct {
	Name  string `json:"name"`
	Login string `json:"login"`
}
