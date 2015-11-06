package document

/*
A Document represents an instance of a text document and a series of
pattern matches for a specified search term.
*/
type Document struct {
	DocumentName string          `json:"documentName"`
	Matches      []*PatternMatch `json:"matches"`
}

/*
HasMatchIndex returns true/false if this document has a match captured
at this location/index in the text document
*/
func (doc *Document) HasMatchIndex(indexToFind int) bool {
	for _, match := range doc.Matches {
		if indexToFind == match.Location {
			return true
		}
	}

	return false
}

/*
NewDocument creates a new document with a document name and blank
matches
*/
func NewDocument(documentName string) *Document {
	return &Document{
		DocumentName: documentName,
		Matches:      make([]*PatternMatch, 0),
	}
}
