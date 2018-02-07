package gorouter

import (
	"net/http"
	"strings"
)

type (
	Tree struct {
		root *Node
		size int
	}

	Node struct {
		Key      string
		Path     string
		Handle   http.HandlerFunc
		depth    int
		children map[string]*Node
		isLeaf   bool
	}
)

func NewNode(key string, depth int) *Node {
	return &Node{
		Key:      key,
		depth:    depth,
		children: make(map[string]*Node),
	}
}

func NewTree() *Tree {
	return &Tree{
		root: NewNode("/", 1),
		size: 1,
	}
}

func (tree *Tree) Add(pattern string, handle http.HandlerFunc) {
	var parent = tree.root

	if pattern != parent.Key {

		pattern = trimPathPrefix(pattern)
		res := splitPattern(pattern)

		for _, key := range res {
			node, ok := parent.children[key]

			if !ok {
				node = NewNode(key, parent.depth+1)

				parent.children[key] = node
			}

			parent = node
		}

	}

	parent.Handle = handle
	parent.isLeaf = true
	parent.Path = pattern
}

func (tree *Tree) Find(pattern string, isRegex int) (nodes []*Node) {
	var (
		node  = tree.root
		queue []*Node
	)

	if pattern == node.Path {
		nodes = append(nodes, node)
		return
	}

	if isRegex == 0 {
		pattern = trimPathPrefix(pattern)
	}

	res := splitPattern(pattern)
	for _, key := range res {

		child, ok := node.children[key]
		if !ok {
			return
		}

		if pattern == child.Path && isRegex == 0 {
			nodes = append(nodes, child)
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

func trimPathPrefix(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

func splitPattern(pattern string) []string {
	return strings.Split(pattern, "/")
}
