package github

import (
	"log"

	"github.com/google/go-github/v33/github"
)

func shouldRetryError(err error) bool {
	if _, ok := err.(*github.RateLimitError); ok {
		log.Printf("[WARN] Received Rate Limit Error")
		return true
	}
	return false
}
