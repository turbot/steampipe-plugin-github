package models

type Label struct {
	NodeId      string `graphql:"nodeId: id" json:"node_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	Color       string `json:"color"`
}
