package models

import "time"

type RepositoryInteractionAbility struct {
	ExpiresAt                        time.Time
	RepositoryInteractionLimit       string
	RepositoryInteractionLimitOrigin string
}
