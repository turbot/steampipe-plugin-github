package models

import "github.com/shurcooL/githubv4"

// PageInfo returns EndCursor and HasNextPage to facilitate paging.
type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}
