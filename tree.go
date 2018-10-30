package gorouter

import (
	"net/http"
	"strings"
)

type (
	// Tree records node
	Tree struct {
		root *Node
	}

	// Node records any URL params, and executes an end handler.
	Node struct {
		key string
		// path records a request path
		path   string
		handle http.HandlerFunc
		// depth records Node's depth
		depth int
		// children records Node's children node
		children map[string]*Node
		// isPattern flag
		isPattern bool
		// middleware records middleware stack
		middleware []MiddlewareType
	}
)

// NewNode returns a newly initialized Node object that implements the Node
func NewNode(key string, depth int) *Node {
	return &Node{
		key:      key,
		depth:    depth,
		children: make(map[string]*Node),
	}
}

// NewTree returns a newly initialized Tree object that implements the Tree
func NewTree() *Tree {
	return &Tree{
		root: NewNode("/", 1),
	}
}

// Add use `pattern` 、handle、middleware stack as node register to tree
func (tree *Tree) Add(pattern string, handle http.HandlerFunc, middleware ...MiddlewareType) {
	var parent = tree.root

	if pattern != parent.key {

		pattern = trimPathPrefix(pattern)
		res := splitPattern(pattern)

		for _, key := range res {
			node, ok := parent.children[key]

			if !ok {
				node = NewNode(key, parent.depth+1)
				if len(middleware) > 0 {
					node.middleware = append(node.middleware, middleware...)
				}

				parent.children[key] = node
			}

			parent = node
		}

	}
	if len(middleware) > 0 && parent.depth == 1 {
		parent.middleware = append(parent.middleware, middleware...)
	}

	parent.handle = handle
	parent.isPattern = true
	parent.path = pattern
}

// Find returns nodes that the request match the route pattern
func (tree *Tree) Find(pattern string, isRegex int) (nodes []*Node) {
	var (
		node  = tree.root
		queue []*Node
	)

	if pattern == node.path {
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

		if pattern == child.path && isRegex == 0 {
			nodes = append(nodes, child)
			return
		}
		node = child
	}

	queue = append(queue, node)

	for len(queue) > 0 {
		var queueTemp []*Node
		for _, n := range queue {
			if n.isPattern {
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

// trimPathPrefix is short for strings.TrimPrefix with param prefix `/`
func trimPathPrefix(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

// splitPattern is short for strings.Split with param seq `/`
func splitPattern(pattern string) []string {
	return strings.Split(pattern, "/")
}
