package models

type PullRequestTemplate struct {
	Filename string `json:"filename"`
	Body     string `json:"body"`
}
