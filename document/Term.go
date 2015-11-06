package document

import (
	"encoding/json"
	"strings"
)

/*
Term is a structure that represents an instance of a matched pattern
in one or more text documents. This is used as the value of tree nodes,
and is what is returned to callers when searching for terms.
*/
type Term struct {
	Key       string      `json:"key"`
	Documents []*Document `json:"documents"`
}

/*
Compare returns 0, 1, or -1 if the second instance of Term
is the same as, greater than, or less than the first instance.
*/
func (term *Term) Compare(compareToTerm *Term) int {
	return strings.Compare(strings.ToLower(term.Key), strings.ToLower(compareToTerm.Key))
}

/*
Equal returns true/false if two Term keys are the same.
*/
func (term *Term) Equal(compareToTerm *Term) bool {
	return strings.ToLower(term.Key) == strings.ToLower(compareToTerm.Key)
}

/*
NewTerm creates a new Term instance with a key name and blank slice of
documents.
*/
func NewTerm(key string) *Term {
	return &Term{
		Key:       key,
		Documents: make([]*Document, 0),
	}
}

/*
ToJSON returns a string of pretty-print JSON representing
this term.
*/
func (term *Term) ToJSON() string {
	bytes, _ := json.MarshalIndent(term, "", "   ")
	return string(bytes)
}
