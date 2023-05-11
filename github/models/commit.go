package models

import (
	"time"
)

// BasicCommit returns the core fields of a Commit.
type BasicCommit struct {
	Sha           string `graphql:"sha: oid"`
	ShortSha      string `graphql:"shortSha: abbreviatedOid"`
	AuthoredDate  time.Time
	Author        GitActor
	CommittedDate time.Time
	Committer     GitActor
	Message       string
	Url           string
}

// Commit returns the full detail of a Commit
type Commit struct {
	BasicCommit
	Additions           int
	AuthoredByCommitter bool
	ChangedFiles        int `graphql:"changedFiles: changedFilesIfAvailable"`
	CommittedViaWeb     bool
	CommitUrl           string
	Deletions           int
	Signature           Signature
	TarballUrl          string
	TreeUrl             string
	CanSubscribe        bool   `graphql:"canSubscribe: viewerCanSubscribe"`
	Subscription        string `graphql:"subscription: viewerSubscription"`
	ZipballUrl          string
	MessageHeadline     string
	Status              CommitStatus
	// AssociatedPullRequests [Pageable]
	// Authors [Pageable]
	// Blame [n-level nesting for an array, requires a path, etc]
	// CheckSuites [Pageable]
	// Comments [Pageable]
	// Deployments [Pageable]
	// File [Requires Path]
	// History [Pageable]
	// OnBehalfOf Organization
	// Parents [Pageable]
}

type CommitStatus struct {
	State string
}
