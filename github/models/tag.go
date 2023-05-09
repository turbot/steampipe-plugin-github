package models

import "time"

type TagWithCommits struct {
	Name   string
	Target struct {
		Commit Commit `graphql:"... on Commit"`
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
				Commit Commit `graphql:"... on Commit"`
			}
		} `graphql:"... on Tag"`
	}
}
