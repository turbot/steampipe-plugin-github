package models

type License struct {
	BasicLicense
	Body           string        `json:"body"`
	Conditions     []LicenseRule `json:"conditions"`
	Description    string        `json:"description"`
	Featured       bool          `json:"featured"`
	Hidden         bool          `json:"hidden"`
	Implementation string        `json:"implementation"`
	Limitations    []LicenseRule `json:"limitations"`
	Permissions    []LicenseRule `json:"permissions"`
	PseudoLicense  bool          `json:"pseudo_license"`
}

type BasicLicense struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	SpdxId   string `json:"spdx_id"`
	Url      string `json:"url"`
}

type LicenseRule struct {
	Description string
	Key         string
	Label       string
}
