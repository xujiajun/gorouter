package gorouter

import (
	"net/http"
)

type (
	Tree struct {
		root *Node
		size int
	}

	Node struct {
		Char     rune
		Path     string
		Handle   http.HandlerFunc
		depth    int
		children map[rune]*Node
		isLeaf   bool
	}
)

func NewNode(char rune, depth int) *Node {
	return &Node{
		Char:     char,
		depth:    depth,
		children: make(map[rune]*Node),
	}
}

func NewTree() *Tree {
	return &Tree{
		root: NewNode(' ', 1),
		size: 1,
	}
}

func (tree *Tree) Add(key string, handle http.HandlerFunc) {
	var parent = tree.root
	allChars := []rune(key)
	for _, char := range allChars {
		node, ok := parent.children[char]

		if !ok {
			node = NewNode(char, parent.depth+1)
			parent.children[char] = node
		}

		parent = node
	}

	parent.Handle = handle
	parent.isLeaf = true
	parent.Path = key
}

func (tree *Tree) Find(key string) (nodes []*Node) {
	var (
		node  = tree.root
		queue []*Node
	)

	allChars := []rune(key)

	for _, char := range allChars {

		child, ok := node.children[char]
		if !ok {
			return
		}
		node = child

	}

	queue = append(queue, node)

	for len(queue) > 0 {
		var queueTemp []*Node
		for _, n := range queue {
			if n.isLeaf {
				nodes = append(nodes, n)
			}

			for _, vnode := range n.children {
				queueTemp = append(queueTemp, vnode)
			}
		}

		queue = queueTemp
	}

	return
}
