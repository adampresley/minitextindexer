package config

import "regexp"

/*
A TextPattern is a structure that describes a regular expression for
capturing text in documents. The Key tells which capture from the
regex capture groups that should be used as the key for tree nodes.
*/
type TextPattern struct {
	Key     int    `json:"key"`
	Pattern string `json:"pattern"`
	Regex   *regexp.Regexp
}
