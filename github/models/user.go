package models

// basicIdentifiers is used to store basic identifying information.
type basicIdentifiers struct {
	Id     int    `graphql:"id: databaseId" json:"id,omitempty"`
	NodeId string `graphql:"nodeId: id" json:"node_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// TODO: Build out user struct
type User struct {
	basicIdentifiers
	Login string `json:"login"`
}
