package models

// BasicCommit returns the core fields of a Commit.
type BasicCommit struct {
	Sha           string       `graphql:"sha: oid" json:"sha"`
	ShortSha      string       `graphql:"shortSha: abbreviatedOid" json:"short_sha"`
	AuthoredDate  NullableTime `json:"authored_date"`
	Author        GitActor     `json:"author"`
	CommittedDate NullableTime `json:"committed_date"`
	Committer     GitActor     `json:"committer"`
	Message       string       `json:"message"`
	Url           string       `json:"url"`
}

// Commit returns the full detail of a Commit
type Commit struct {
	Sha                 string       `graphql:"sha: oid" json:"sha"`
	ShortSha            string       `graphql:"shortSha: abbreviatedOid @include(if:$includeCommitShortSha)" json:"short_sha"`
	AuthoredDate        NullableTime `graphql:"authoredDate @include(if:$includeCommitAuthoredDate)" json:"authored_date"`
	Author              GitActor     `graphql:"author @include(if:$includeCommitAuthor)" json:"author"`
	CommittedDate       NullableTime `graphql:"committedDate @include(if:$includeCommitCommittedDate)" json:"committed_date"`
	Committer           GitActor     `graphql:"committer @include(if:$includeCommitCommitter)" json:"committer"`
	Message             string       `graphql:"message @include(if:$includeCommitMessage)" json:"message"`
	Url                 string       `graphql:"url @include(if:$includeCommitUrl)" json:"url"`
	Additions           int          `graphql:"additions @include(if:$includeCommitAdditions)" json:"additions"`
	AuthoredByCommitter bool         `graphql:"authoredByCommitter @include(if:$includeCommitAuthoredByCommitter)" json:"authored_by_committer"`
	ChangedFiles        int          `graphql:"changedFiles: changedFilesIfAvailable @include(if:$includeCommitChangedFiles)" json:"changed_files"`
	CommittedViaWeb     bool         `graphql:"committedViaWeb @include(if:$includeCommitCommittedViaWeb)" json:"committed_via_web"`
	CommitUrl           string       `graphql:"commitUrl @include(if:$includeCommitCommitUrl)" json:"commit_url"`
	Deletions           int          `graphql:"deletions @include(if:$includeCommitDeletions)" json:"deletions"`
	Signature           Signature    `graphql:"signature @include(if:$includeCommitSignature)" json:"signature"`
	TarballUrl          string       `graphql:"tarballUrl @include(if:$includeCommitTarballUrl)" json:"tarball_url"`
	TreeUrl             string       `graphql:"treeUrl @include(if:$includeCommitTreeUrl)" json:"tree_url"`
	CanSubscribe        bool         `graphql:"canSubscribe: viewerCanSubscribe @include(if:$includeCommitCanSubscribe)" json:"can_subscribe"`
	Subscription        string       `graphql:"subscription: viewerSubscription @include(if:$includeCommitSubscription)" json:"subscription"`
	ZipballUrl          string       `graphql:"zipballUrl @include(if:$includeCommitZipballUrl)" json:"zipball_url"`
	MessageHeadline     string       `graphql:"messageHeadline @include(if:$includeCommitMessageHeadline)" json:"message_headline"`
	Status              CommitStatus `graphql:"status @include(if:$includeCommitStatus)" json:"status"`
	NodeId              string       `graphql:"nodeId:id @include(if:$includeCommitNodeId)" json:"node_id"`
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

type BaseCommit struct {
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
	NodeId              string       `graphql:"nodeId:id" json:"node_id"`
}
