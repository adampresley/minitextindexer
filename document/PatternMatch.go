package document

/*
A PatternMatch is the information about a particlar match in a document.
This includes the starting location/index of the match, the contents of the
match, and all regex capture groups.
*/
type PatternMatch struct {
	Captures []string `json:"captures"`
	Location int      `json:"location"`
	Match    string   `json:"match"`
}
