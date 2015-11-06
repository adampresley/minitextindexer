package document

/*
A FileIndexMatch is used when scanning a file for text patterns
*/
type FileIndexMatch struct {
	Captures []string
	Key      string
	Location int
	Match    string
}
