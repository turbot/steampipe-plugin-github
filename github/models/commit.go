package models

import (
	"time"
)

// BasicCommit returns the core fields of a Commit.
type BasicCommit struct {
	Sha          string `graphql:"sha: oid"`
	ShortSha     string `graphql:"shortSha: abbreviatedOid"`
	AuthoredDate time.Time
	Message      string
	Author       GitActor
	Url          string
}

// Commit returns the full detail of a Commit
type Commit struct {
	BasicCommit
	Additions           int
	AuthoredByCommitter bool
	ChangedFiles        int `graphql:"changedFiles: changedFilesIfAvailable"`
	CommittedDate       time.Time
	CommittedViaWeb     bool
	Committer           GitActor
	CommitUrl           string
	Deletions           int
	Signature           Signature
	TarballUrl          string
	TreeUrl             string
	CanSubscribe        bool   `graphql:"canSubscribe: viewerCanSubscribe"`
	Subscription        string `graphql:"subscription: viewerSubscription"`
	ZipballUrl          string
	// AssociatedPullRequests
	// Authors
	// Blame
	// CheckSuites
	// Comments
	// Deployments
	// File
	// History
	// OnBehalfOf Organization
	// Parents
	// Status
}
