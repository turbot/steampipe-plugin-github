package models

import "time"

// RateLimit information
type BaseRateLimit struct {
	Remaining int       `graphql:"remaining @include(if:$includeRLRemaining)" json:"remaining"`
	Used      int       `graphql:"used @include(if:$includeRLUsed)" json:"used"`
	Cost      int       `graphql:"cost @include(if:$includeRLCost)" json:"cost"`
	Limit     int       `graphql:"limit @include(if:$includeRLLimit)" json:"limit"`
	ResetAt   time.Time `graphql:"resetAt @include(if:$includeRLResetAt)" json:"reset_at"`
	NodeCount int       `graphql:"remaining @include(if:$includeRLNodeCount)" json:"node_count"`
}

type RateLimit struct {
	Remaining int
	Used      int
	Cost      int
	Limit     int
	ResetAt   time.Time
	NodeCount int
}
