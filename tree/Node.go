package tree

import "github.com/adampresley/minitextindexer/document"

/*
Node is an individual element in the tree
*/
type Node struct {
	Value *document.Term `json:"value"`

	Left   *Node `json:"left"`
	Right  *Node `json:"right"`
	Parent *Node `json:"-"`
}

/*
FindDocument searches this node's value for a specific document.
It returns nil if one is not found.
*/
func (node *Node) FindDocument(documentName string) *document.Document {
	for _, document := range node.Value.Documents {
		if document.DocumentName == documentName {
			return document
		}
	}

	return nil
}

/*
NewNode creates a new tree node with a specified value.
*/
func NewNode(term *document.Term) *Node {
	return &Node{
		Value: term,
	}
}
