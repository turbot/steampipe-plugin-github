package models

type TreeEntry struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	LineCount   int      `json:"line_count"`
	Size        int      `json:"size"`
	IsGenerated bool     `json:"is_generated"`
	Language    Language `json:"language"`
	Extension   string   `json:"extension"`
	Mode        int      `json:"mode"`
	Object      struct {
		Blob Blob `graphql:"... on Blob" json:"blob"`
	} `json:"object"`
}

type Blob struct {
	NodeId      string `graphql:"nodeId: id" json:"node_id"`
	IsTruncated bool   `json:"is_truncated"`
	IsBinary    bool   `json:"is_binary"`
	Text        string `json:"text"`
	ByteSize    int    `json:"byte_size"`
	CommitSha   string `graphql:"commitSha: oid" json:"commit_sha"`
	CommitUrl   string `json:"commit_url"`
}
