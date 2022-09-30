package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func shouldRetryError(ctx context.Context, err error) bool {
	if _, ok := err.(*github.RateLimitError); ok {
		plugin.Logger(ctx).Debug("errors.shouldRetryError", "rate_limit_error", err, "rate", err.(*github.RateLimitError).Rate)
		return false
	}

	if _, ok := err.(*github.AbuseRateLimitError); ok {
		var retryAfter *time.Duration
		if err.(*github.AbuseRateLimitError).RetryAfter != nil {
			retryAfter = err.(*github.AbuseRateLimitError).RetryAfter
		}
		plugin.Logger(ctx).Debug("errors.shouldRetryError", "abuse_rate_limit_error", err, "retryAfter", retryAfter)
		return true
	}

	if _, ok := err.(*github.AbuseRateLimitError); ok {
		log.Printf("[WARN] Received Secondary Rate Limit Error")
		return true
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
