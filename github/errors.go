package github

import (
	"log"
	"strconv"

	"github.com/google/go-github/v32/github"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// function which returns an ErrorPredicate for Not Found Errors
func shouldIgnoreError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if githubErr, ok := err.(*github.ErrorResponse); ok {
			return helpers.StringSliceContains(notFoundErrors, strconv.Itoa(githubErr.Response.StatusCode))
		}
		return false
	}
}

func shouldRetryError(err error) bool {
	if _, ok := err.(*github.RateLimitError); ok {
		log.Printf("[WARN] Received Rate Limit Error")
		return true
	}
	return false
}
