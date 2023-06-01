package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func shouldRetryError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
	if _, ok := err.(*github.AbuseRateLimitError); ok {
		var retryAfter *time.Duration
		if err.(*github.AbuseRateLimitError).RetryAfter != nil {
			retryAfter = err.(*github.AbuseRateLimitError).RetryAfter
		}
		plugin.Logger(ctx).Debug("github_errors.shouldRetryError", "abuse_rate_limit_error", err, "retry_after", retryAfter)
		return true
	}

	if _, ok := err.(*github.RateLimitError); ok {
		// Get the limit reset timestamp if returned
		var resetAfter time.Time
		if err.(*github.RateLimitError).Rate.String() != "" {
			resetAfter = err.(*github.RateLimitError).Rate.Reset.Time
		}

		// Get the remaining time
		t1 := time.Now()
		diff := resetAfter.Sub(t1).Seconds()
		plugin.Logger(ctx).Debug("github_errors.shouldRetryError", "rate_limit_error", err, "reset_after", diff)

		// Treat the error as non-fatal if the remaining time for limit reset is
		// less than 60s
		return diff <= 60
	}

	// v4 secondary rate limit
	if strings.Contains(err.Error(), "You have exceeded a secondary rate limit.") {
		plugin.Logger(ctx).Debug("github_errors.shouldRetryError", "abuse_rate_limit_error", err)
		return true
	}

	return false
}

func retryConfig() *plugin.RetryConfig {
	return &plugin.RetryConfig{
		ShouldRetryErrorFunc: shouldRetryError,
		MaxAttempts:          10,
		BackoffAlgorithm:     "Exponential",
		RetryInterval:        1000,
		CappedDuration:       30000,
	}
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
