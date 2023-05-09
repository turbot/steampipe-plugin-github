package models

// Signature returns information about signatures including the Email, State and validity (IsValid).
type Signature struct {
	Email             string
	IsValid           bool
	State             string
	WasSignedByGitHub bool
	Signer            struct {
		Email string
		Login string
	}
}
