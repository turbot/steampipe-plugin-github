package models

import "github.com/shurcooL/githubv4"

type Milestone struct {
	Closed             bool                    `json:"closed"`
	ClosedAt           NullableTime            `json:"closed_at"`
	CreatedAt          NullableTime            `json:"created_at"`
	Creator            Actor                   `json:"creator"`
	Description        string                  `json:"description"`
	DueOn              NullableTime            `json:"due_on"`
	Number             int                     `json:"number"`
	ProgressPercentage float32                 `json:"progress_percentage"`
	State              githubv4.MilestoneState `json:"state"`
	Title              string                  `json:"title"`
	UpdatedAt          NullableTime            `json:"updated_at"`
	UserCanClose       bool                    `graphql:"userCanClose: viewerCanClose" json:"user_can_close"`
	UserCanReopen      bool                    `graphql:"userCanReopen: viewerCanReopen" json:"user_can_reopen"`
}
