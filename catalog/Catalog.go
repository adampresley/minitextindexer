package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/adampresley/minitextindexer/config"
	"github.com/adampresley/minitextindexer/document"
	"github.com/adampresley/minitextindexer/tree"

	"github.com/adampresley/directorywatcher"
	"github.com/adampresley/logging"
)

/*
A Catalog represents a physical file tree and its indexed, virtual tree.
*/
type Catalog struct {
	sync.Mutex

	basePaths    []string
	config       *config.Configuration
	log          *logging.Logger
	textPatterns []*config.TextPattern
	tree         *tree.Tree
	watchers     []*directorywatcher.DirectoryWatcher
}

func (catalog *Catalog) compileRegexes() {
	for index, textPattern := range catalog.textPatterns {
		exp, err := regexp.Compile(textPattern.Pattern)
		if err != nil {
			catalog.log.Errorf("Problem compiling regex [%s]: %s", textPattern.Pattern, err.Error())
		} else {
			catalog.textPatterns[index].Regex = exp
		}
	}
}

/*
FindTerm searches the tree for a specific term.
*/
func (catalog *Catalog) FindTerm(searchTerm string) *document.Term {
	node := catalog.tree.Find(document.NewTerm(searchTerm))

	if node == nil {
		return nil
	}

	return node.Value
}

/*
Index creates the virtual index tree. This operation locks the
catalog.
*/
func (catalog *Catalog) Index() error {
	startTime := time.Now()
	fileCount := 0
	nodeCount := 0

	catalog.Lock()

	doneChannel := make(chan bool)
	indexChannel := make(chan document.DocumentIndex, 100)
	waitGroup := &sync.WaitGroup{}

	go func(waitGroup *sync.WaitGroup) {
		for {
			select {
			case indexItem := <-indexChannel:
				for key, newDocument := range indexItem {
					/*
					 * Create a new Term and add the document to it.
					 */
					termToFind := document.NewTerm(key)
					termToFind.Documents = append(termToFind.Documents, newDocument)

					/*
					 * See if this term is already in the tree. If not, just add it.
					 * If it is there, we need to see if we have this document captured
					 * already. If not, add it. If so, only add our matches if we don't
					 * already have those too.
					 */
					existingTermNode := catalog.tree.Find(termToFind)

					if existingTermNode == nil {
						catalog.tree.Add(termToFind)
						nodeCount++
					} else {
						existingDocument := existingTermNode.FindDocument(newDocument.DocumentName)

						if existingDocument == nil {
							existingTermNode.Value.Documents = append(existingTermNode.Value.Documents, newDocument)
						} else {
							currentMatches := existingDocument.Matches

							for _, newMatch := range newDocument.Matches {
								if !existingDocument.HasMatchIndex(newMatch.Location) {
									currentMatches = append(currentMatches, newMatch)
								}
							}

							existingDocument.Matches = currentMatches
						}
					}
				}

				waitGroup.Done()

			case <-doneChannel:
				catalog.log.Debug("Done indexing catalog")
				break
			}
		}
	}(waitGroup)

	for _, basePath := range catalog.config.Paths {
		filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				/*
				 * Only index this file if it matches the configured file pattern
				 */
				isFilePatternMatch := false

				for _, filePattern := range catalog.config.FilePatterns {
					if strings.Contains(path, filePattern) {
						isFilePatternMatch = true
						break
					}
				}

				if !isFilePatternMatch {
					return nil
				}

				/*
				 * Index this file
				 */
				fileCount++
				waitGroup.Add(1)

				file := document.NewPhysicalFile(path, catalog.textPatterns)
				_, err := file.Read()
				if err != nil {
					catalog.log.Errorf("Error reading file %s: %s", path, err.Error())
					return err
				}

				index := file.CreateIndex()
				indexChannel <- index
			}

			return nil
		})
	}

	waitGroup.Wait()
	doneChannel <- true
	catalog.Unlock()

	catalog.log.Infof("Time to index %d files with %d nodes: %s", fileCount, nodeCount, time.Since(startTime))
	return nil
}

/*
NewCatalog returns a new instance of a Catalog structure. It will
create the initial index and start a directory watcher for the physical
file structure.
*/
func NewCatalog(log *logging.Logger, config *config.Configuration) *Catalog {
	// TODO: Using "mn" as a root node. This may need to be calculated based on found elements.
	// If too many files come on side or the other of "mn" we will have an unbalanced tree.
	catalog := &Catalog{
		basePaths:    config.Paths,
		config:       config,
		log:          log,
		textPatterns: config.TextPatterns,
		tree:         tree.NewTree(document.NewTerm("mn")),
	}

	/*
	 * Define a directory watcher function to be used by each directory watcher
	 */
	watcherFunc := func(path string, info os.FileInfo, startTime time.Time, modificationTime time.Time) error {
		isFilePatternMatch := false

		for _, filePattern := range config.FilePatterns {
			if strings.Contains(path, filePattern) {
				isFilePatternMatch = true
				break
			}
		}

		if isFilePatternMatch {
			log.Infof("Detected change in path %s", path)
		}

		return nil
	}

	catalog.watchers = make([]*directorywatcher.DirectoryWatcher, len(config.Paths))

	for index, basePath := range config.Paths {
		catalog.watchers[index] = directorywatcher.NewDirectoryWatcher(basePath, log)
		catalog.watchers[index].Watch(watcherFunc)
	}

	catalog.compileRegexes()
	return catalog
}

/*
Search searches the tree for nodes containing a term. This uses a depth-first
pre-order traversal.
*/
func (catalog *Catalog) Search(searchTerm string) []*document.Term {
	nodes := catalog.tree.Search(searchTerm)

	if nodes == nil {
		return nil
	}

	results := make([]*document.Term, len(nodes))

	for index, node := range nodes {
		results[index] = node.Value
	}

	return results
}

/*
ToJSON returns a pretty printed string of this catalog as JSON
*/
func (catalog *Catalog) ToJSON() string {
	result := make(map[string]interface{})
	catalog.Lock()

	result["tree"] = catalog.tree
	result["basePaths"] = catalog.basePaths

	bytes, _ := json.MarshalIndent(result, "", "   ")

	catalog.Unlock()
	return string(bytes)
}
