package tree

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/adampresley/minitextindexer/document"
)

/*
Tree is a tree structure which is the index of our
text files and their matching terms/values
*/
type Tree struct {
	Root *Node `json:"root"`
}

/*
Add creates a new tree node and inserts it into the tree.
*/
func (tree *Tree) Add(term *document.Term) *Node {
	newNode := NewNode(term)
	parent := tree.findLast(term)

	if parent == nil {
		tree.Root = newNode
	} else {
		compare := term.Compare(parent.Value)

		if compare < 0 {
			parent.Left = newNode
		} else if compare > 0 {
			parent.Right = newNode
		} else {
			// value is already in the tree
			return nil
		}

		newNode.Parent = parent
	}

	return newNode
}

/*
Find searches for a specific term in the tree. If it is not found
nil is returned.
*/
func (tree *Tree) Find(term *document.Term) *Node {
	currentNode := tree.Root

	for {
		if currentNode == nil {
			break
		}

		compare := term.Compare(currentNode.Value)

		if compare < 0 {
			currentNode = currentNode.Left
		} else if compare > 0 {
			currentNode = currentNode.Right
		} else {
			return currentNode
		}
	}

	return nil
}

/*
Traverse the tree until we find a suitable node that has a lexical
ordering just before the specified term.
*/
func (tree *Tree) findLast(term *document.Term) *Node {
	currentNode := tree.Root
	var previousNode *Node

	for {
		if currentNode == nil {
			break
		}

		previousNode = currentNode
		compare := term.Compare(currentNode.Value)

		if compare < 0 {
			currentNode = currentNode.Left
		} else if compare > 0 {
			currentNode = currentNode.Right
		} else {
			return currentNode
		}
	}

	return previousNode
}

func inOrder(node *Node, nodeChannel chan *Node, searchTerm string, waitGroup *sync.WaitGroup) {
	if node == nil {
		return
	}

	waitGroup.Add(1)
	nodeChannel <- node
	inOrder(node.Left, nodeChannel, searchTerm, waitGroup)
	inOrder(node.Right, nodeChannel, searchTerm, waitGroup)
}

/*
NewTree creates a new tree with a specific root node
*/
func NewTree(rootTerm *document.Term) *Tree {
	return &Tree{
		Root: NewNode(rootTerm),
	}
}

/*
Search returns a set of nodes who's values contain a search term
*/
func (tree *Tree) Search(searchTerm string) []*Node {
	nodeChannel := make(chan *Node, 100)
	doneChannel := make(chan bool)
	waitGroup := &sync.WaitGroup{}
	var results []*Node

	go func(nodeChannel chan *Node) {
		for {
			select {
			case node := <-nodeChannel:
				if strings.Contains(strings.ToLower(node.Value.Key), strings.ToLower(searchTerm)) {
					results = append(results, node)
				}

				waitGroup.Done()

			case <-doneChannel:
				break
			}
		}
	}(nodeChannel)

	inOrder(tree.Root, nodeChannel, searchTerm, waitGroup)
	waitGroup.Wait()
	doneChannel <- true

	return results
}

/*
ToJSON returns a pretty printed string of this tree as JSON
*/
func (tree *Tree) ToJSON() string {
	bytes, _ := json.MarshalIndent(tree, "", "   ")
	return string(bytes)
}
