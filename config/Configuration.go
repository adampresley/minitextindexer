package config

/*
A Configuration structure represents the data necessary to configure
a Mini Text Indexer instance.
*/
type Configuration struct {
	FilePatterns []string       `json:"filePatterns"`
	Paths        []string       `json:"paths"`
	TextPatterns []*TextPattern `json:"textPatterns"`
}
