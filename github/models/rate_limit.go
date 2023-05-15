package models

import "time"

// RateLimit information
type RateLimit struct {
	Remaining int
	Used      int
	Cost      int
	Limit     int
	ResetAt   time.Time
	NodeCount int
}
