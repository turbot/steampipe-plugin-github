package models

import (
	"time"
)

// BasicCommit returns the core fields of a Commit.
type BasicCommit struct {
	Sha           string    `graphql:"sha: oid" json:"sha"`
	ShortSha      string    `graphql:"shortSha: abbreviatedOid" json:"short_sha"`
	AuthoredDate  time.Time `json:"authored_date"`
	Author        GitActor  `json:"author"`
	CommittedDate time.Time `json:"committed_date"`
	Committer     GitActor  `json:"committer"`
	Message       string    `json:"message"`
	Url           string    `json:"url"`
}

// Commit returns the full detail of a Commit
type Commit struct {
	BasicCommit
	Additions           int          `json:"additions"`
	AuthoredByCommitter bool         `json:"authored_by_committer"`
	ChangedFiles        int          `graphql:"changedFiles: changedFilesIfAvailable" json:"changed_files"`
	CommittedViaWeb     bool         `json:"committed_via_web"`
	CommitUrl           string       `json:"commit_url"`
	Deletions           int          `json:"deletions"`
	Signature           Signature    `json:"signature"`
	TarballUrl          string       `json:"tarball_url"`
	TreeUrl             string       `json:"tree_url"`
	CanSubscribe        bool         `graphql:"canSubscribe: viewerCanSubscribe" json:"can_subscribe"`
	Subscription        string       `graphql:"subscription: viewerSubscription" json:"subscription"`
	ZipballUrl          string       `json:"zipball_url"`
	MessageHeadline     string       `json:"message_headline"`
	Status              CommitStatus `json:"status"`
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
	State string `json:"state"`
}
