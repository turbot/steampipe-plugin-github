package models

type License struct {
	BasicLicense
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
	Name     string `graphql:"name @include(if:$includeLicenseName)" json:"name"`
	Nickname string `graphql:"nickname @include(if:$includeLicenseNickname)" json:"nickname"`
	SpdxId   string `graphql:"spdxId @include(if:$includeLicenseSpdxId)" json:"spdx_id"`
	Url      string `graphql:"url @include(if:$includeLicenseUrl)" json:"url"`
}

type LicenseRule struct {
	Description string
	Key         string
	Label       string
}
