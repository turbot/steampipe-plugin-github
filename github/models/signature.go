package models

// Signature returns information about signatures including the Email, State and validity (IsValid).
type Signature struct {
	Email             string `json:"email"`
	IsValid           bool   `json:"is_valid"`
	State             string `json:"state"`
	WasSignedByGitHub bool   `json:"was_signed_by_github"`
	Signer            struct {
		Email string `json:"email"`
		Login string `json:"login"`
	} `json:"signer"`
}
