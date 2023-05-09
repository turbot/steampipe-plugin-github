package models

// basicIdentifiers is used to store basic identifying information.
type basicIdentifiers struct {
	Id     int    `graphql:"id: databaseId"`
	NodeId string `graphql:"nodeId: id"`
	Name   string
}

// TODO: Build out user struct
type User struct {
	basicIdentifiers
	Login string
}
