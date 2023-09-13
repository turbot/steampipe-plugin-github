package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func shouldRetryError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
	if e, ok := err.(*github.AbuseRateLimitError); ok {
		var retryAfter *time.Duration
		if e.RetryAfter != nil {
			retryAfter = e.RetryAfter
		}
		plugin.Logger(ctx).Debug("github_errors.shouldRetryError", "abuse_rate_limit_error", err, "retry_after", retryAfter)
		return true
	}

	if e, ok := err.(*github.RateLimitError); ok {
		// Get the limit reset timestamp if returned
		var resetAfter time.Time
		if e.Rate.String() != "" {
			resetAfter = e.Rate.Reset.Time
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

	// v4 execution timeout error
	if strings.Contains(err.Error(), "Something went wrong while executing your query. This may be the result of a timeout, or it could be a GitHub bug.") {
		plugin.Logger(ctx).Debug("github_errors.shouldRetryError", "execution_timeout_error", err)
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

// function which returns an ErrorPredicate for GitHub API calls
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
