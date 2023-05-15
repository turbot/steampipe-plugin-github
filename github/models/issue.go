package models

type IssueTemplate struct {
	About    string `json:"about"`
	Body     string `json:"body"`
	Filename string `json:"filename"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	// Assignees [pageable]
	// Labels [pageable]
}
