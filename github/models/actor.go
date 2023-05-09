package models

import "time"

type Actor struct {
	AvatarUrl string
	Login     string
	Url       string
}

type GitActor struct {
	AvatarUrl string
	Date      time.Time
	Email     string
	Name      string
	User      User
}
