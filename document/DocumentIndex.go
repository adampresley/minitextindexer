package document

/*
A DocumentIndex represents a matched pattern text in a single document. This
is used during the document scanning process.
*/
type DocumentIndex map[string]*Document
