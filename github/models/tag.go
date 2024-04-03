package models

import "time"

type TagWithCommits struct {
	Name   string `graphql:"name @include(if:$includeTagName)" json:"name"`
	Target struct {
		Commit BaseCommit `graphql:"... on Commit"`
		Tag    struct {
			Message string
			Tagger  struct {
				Name string
				Date time.Time
				User struct {
					Login string
				}
			}
			Target struct {
				Commit BaseCommit `graphql:"... on Commit"`
			}
		} `graphql:"... on Tag"`
	} `graphql:"target @include(if:$includeTagTarget)" json:"target"`
}
