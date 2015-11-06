package document

import (
	"io/ioutil"
	"sync"

	"github.com/adampresley/minitextindexer/config"
)

/*
A PhysicalFile represents a file on a file system. This structure
provides methods to read a file and perform an pattern match
scan on it.
*/
type PhysicalFile struct {
	Contents string `json:"contents"`
	FileName string `json:"fileName"`

	textPatterns []*config.TextPattern
}

/*
CreateIndex creates an index of matched patterns, each with a set
of documents attached to them.
*/
func (file *PhysicalFile) CreateIndex() DocumentIndex {
	result := make(DocumentIndex)
	matchChannel := make(chan FileIndexMatch, 5)
	doneChan := make(chan bool)

	waitGroup := &sync.WaitGroup{}

	/*
	 * Goroutine to read location indexes and terms off a channel and put them
	 * into our slices
	 */
	go func(waitGroup *sync.WaitGroup) {
		for {
			select {
			case match := <-matchChannel:
				newPatternMatch := &PatternMatch{
					Captures: match.Captures,
					Location: match.Location,
					Match:    match.Match,
				}

				if document, ok := result[match.Key]; ok {
					document.Matches = append(document.Matches, newPatternMatch)
					result[match.Key] = document
				} else {
					newDocument := NewDocument(file.FileName)
					newDocument.Matches = append(newDocument.Matches, newPatternMatch)
					result[match.Key] = newDocument
				}

				waitGroup.Done()

			case <-doneChan:
				break
			}
		}
	}(waitGroup)

	/*
	 * Iterate over each keyword and fire off a goroutine to search
	 * this file's contents for the keyword. When a keyword is found
	 * write it to the term channel, and the match locations get
	 * written to the locations channel.
	 */
	for _, textPattern := range file.textPatterns {
		searchResult := textPattern.Regex.FindAllStringSubmatch(file.Contents, -1)
		searchResultIndexes := textPattern.Regex.FindAllStringIndex(file.Contents, -1)

		if len(searchResult) > 0 {
			waitGroup.Add(len(searchResult))

			for matchIndex, matchedSet := range searchResult {
				match := FileIndexMatch{
					Captures: matchedSet,
					Key:      matchedSet[textPattern.Key],
					Location: searchResultIndexes[matchIndex][0],
					Match:    matchedSet[0],
				}

				matchChannel <- match
			}
		}
	}

	waitGroup.Wait()
	doneChan <- true

	return result
}

/*
NewPhysicalFile creates a new PhysicalFile structure. It takes a file name
and a set of regular expressions to run against it.
*/
func NewPhysicalFile(fileName string, textPatterns []*config.TextPattern) *PhysicalFile {
	return &PhysicalFile{
		FileName:     fileName,
		textPatterns: textPatterns,
	}
}

/*
Read reads the actual contents of a file and puts them into
the Contents key
*/
func (file *PhysicalFile) Read() (string, error) {
	bytes, err := ioutil.ReadFile(file.FileName)
	file.Contents = string(bytes)
	return file.Contents, err
}

/*
String returns the file name
*/
func (file *PhysicalFile) String() string {
	return file.FileName
}
