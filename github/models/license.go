package models

type License struct {
	Body           string
	Conditions     []LicenseRule
	Description    string
	Featured       bool
	Hidden         bool
	Implementation string
	Key            string
	Limitations    []LicenseRule
	Name           string
	Nickname       string
	Permissions    []LicenseRule
	PseudoLicense  bool
	SpdxId         string
	Url            string
}

type LicenseRule struct {
	Description string
	Key         string
	Label       string
}
