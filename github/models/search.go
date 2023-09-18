package models

type TextMatch struct {
	Fragment   string               `json:"fragment"`
	Property   string               `json:"property"`
	Highlights []TextMatchHighlight `json:"highlights"`
}

type TextMatchHighlight struct {
	BeginIndice int    `json:"begin_indice"`
	EndIndice   int    `json:"end_indice"`
	Text        string `json:"text"`
}

type SearchRepositoryResult struct {
	TextMatches []TextMatch
	Node        struct {
		Repository `graphql:"... on Repository"`
	}
}
