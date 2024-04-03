package models

type License struct {
	Key            string        `json:"key"`
	Name           string        `graphql:"name @include(if:$includeLicenseName)" json:"name"`
	Nickname       string        `graphql:"nickname @include(if:$includeLicenseNickname)" json:"nickname"`
	SpdxId         string        `graphql:"spdxId @include(if:$includeLicenseSpdxId)" json:"spdx_id"`
	Url            string        `graphql:"url @include(if:$includeLicenseUrl)" json:"url"`
	Body           string        `json:"body"`
	Conditions     []LicenseRule `graphql:"conditions @include(if:$includeLicenseConditions)" json:"conditions"`
	Description    string        `graphql:"description @include(if:$includeLicenseDescription)" json:"description"`
	Featured       bool          `graphql:"featured @include(if:$includeLicenseFeatured)" json:"featured"`
	Hidden         bool          `graphql:"hidden @include(if:$includeLicenseHidden)" json:"hidden"`
	Implementation string        `graphql:"implementation @include(if:$includeLicenseImplementation)" json:"implementation"`
	Limitations    []LicenseRule `graphql:"limitations @include(if:$includeLicenseLimitations)" json:"limitations"`
	Permissions    []LicenseRule `graphql:"permissions @include(if:$includeLicensePermissions)" json:"permissions"`
	PseudoLicense  bool          `graphql:"pseudoLicense @include(if:$includeLicensePseudoLicense)" json:"pseudo_license"`
}

type BasicLicense struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	SpdxId   string `json:"spdx_id"`
	Url      string `json:"url"`
}

type BaseLicense struct {
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

type LicenseRule struct {
	Description string
	Key         string
	Label       string
}
