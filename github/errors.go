package github

import (
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func shouldRetryError(err error) bool {
	if _, ok := err.(*github.RateLimitError); ok {
		log.Printf(`rate limit %+v`, err)
		time.Sleep(1 * time.Minute)
		return true
	}
	log.Printf(`error %+v`, err)
	return false
}

// function which returns an ErrorPredicate for Github API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if err != nil {
			for _, item := range notFoundErrors {
				if strings.Contains(err.Error(), item) {
					return true
				}
			}
		}
		return false
	}
}
