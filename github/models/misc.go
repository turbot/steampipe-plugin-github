package models

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type NullableTime struct {
	time.Time
}

func (t NullableTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	} else {
		return t.Time.MarshalJSON()
	}
}

type NameSlug struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type NameLogin struct {
	Name  string `json:"name"`
	Login string `json:"login"`
}

type BasicRef struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

type BasicRefWithBranchProtectionRule struct {
	BasicRef
	BranchProtectionRule *BaseBranchProtectionRule `json:"branch_protection_rule,omitempty"`
}

type Language struct {
	Id    githubv4.ID `json:"id"`
	Name  string      `json:"name"`
	Color string      `json:"color"`
}

type SponsorsListing struct {
	Id                         string               `json:"id,omitempty"`
	ActiveGoal                 SponsorsGoal         `json:"active_goal"`
	ActiveStripeConnectAccount StripeConnectAccount `json:"active_stripe_connect_account"`
	BillingCountryOrRegion     string               `json:"billing_country_or_region"`
	ContactEmailAddress        string               `json:"contact_email_address"`
	CreatedAt                  time.Time            `json:"created_at"`
	DashboardUrl               string               `json:"dashboard_url"`
	FullDescription            string               `json:"full_description"`
	IsPublic                   bool                 `json:"is_public"`
	Name                       string               `json:"name"`
	NextPayoutDate             time.Time            `json:"next_payout_date"`
	ResidenceCountryOrRegion   string               `json:"residence_country_or_region"`
	ShortDescription           string               `json:"short_description"`
	Slug                       string               `json:"slug"`
	Url                        string               `json:"url"`
	// Tiers [pageable]
	// FeaturedItems [searchable by key]
}

type SponsorsGoal struct {
	Description     string                    `json:"description"`
	PercentComplete int                       `json:"percent_complete"`
	TargetValue     int                       `json:"target_value"`
	Title           string                    `json:"title"`
	Kind            githubv4.SponsorsGoalKind `json:"kind"`
}

type StripeConnectAccount struct {
	AccountId              string `json:"account_id"`
	BillingCountryOrRegion string `json:"billing_country_or_region"`
	CountryOrRegion        string `json:"country_or_region"`
	IsActive               bool   `json:"is_active"`
	StripeDashboardUrl     string `json:"stripe_dashboard_url"`
}

type Count struct {
	TotalCount int `json:"total_count"`
}
