package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func shouldRetryError(ctx context.Context, err error) bool {
	if _, ok := err.(*github.AbuseRateLimitError); ok {
		var retryAfter *time.Duration
		if err.(*github.AbuseRateLimitError).RetryAfter != nil {
			retryAfter = err.(*github.AbuseRateLimitError).RetryAfter
		}
		plugin.Logger(ctx).Debug("errors.shouldRetryError", "abuse_rate_limit_error", err, "retryAfter", retryAfter)
		return true
	}

	if _, ok := err.(*github.RateLimitError); ok {
		// Get the timestamp after which the limit will get reset
		var resetAfter time.Time
		if err.(*github.RateLimitError).Rate.String() != "" {
			resetAfter = err.(*github.RateLimitError).Rate.Reset.Time
		}

		// Get the remaining time
		t1 := time.Now()
		diff := resetAfter.Sub(t1).Seconds()
		plugin.Logger(ctx).Debug("errors.shouldRetryError", "rate_limit_error", err, "resetAfter", diff)

		// Treat the error as non-fatal if the remaining time for limit reset is less than 60s
		if diff <= 60 {
			return true
		}
	}

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
