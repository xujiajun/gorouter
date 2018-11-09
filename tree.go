package gorouter

import (
	"net/http"
	"strings"
)

type (
	// Tree records node
	Tree struct {
		root       *Node
		parameters Parameters
		routes     map[string]*Node
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
		root:   NewNode("/", 1),
		routes: make(map[string]*Node),
	}
}

// Add use `pattern` 、handle、middleware stack as node register to tree
func (tree *Tree) Add(pattern string, handle http.HandlerFunc, middleware ...MiddlewareType) {
	var currentNode = tree.root

	if pattern != currentNode.key {

		pattern = trimPathPrefix(pattern)
		res := splitPattern(pattern)

		for _, key := range res {
			node, ok := currentNode.children[key]

			if !ok {
				node = NewNode(key, currentNode.depth+1)
				if len(middleware) > 0 {
					node.middleware = append(node.middleware, middleware...)
				}

				currentNode.children[key] = node
			}

			currentNode = node
		}

	}
	if len(middleware) > 0 && currentNode.depth == 1 {
		currentNode.middleware = append(currentNode.middleware, middleware...)
	}

	currentNode.handle = handle
	currentNode.isPattern = true
	currentNode.path = pattern

	if routeName := tree.parameters.routeName; routeName != "" {
		tree.routes[routeName] = currentNode
	}
}

// Find returns nodes that the request match the route pattern
func (tree *Tree) Find(pattern string, isRegex bool) (nodes []*Node) {
	var (
		node  = tree.root
		queue []*Node
	)

	if pattern == node.path {
		nodes = append(nodes, node)
		return
	}

	if !isRegex {
		pattern = trimPathPrefix(pattern)
	}

	res := splitPattern(pattern)

	for _, key := range res {
		child, ok := node.children[key]

		if !ok && isRegex {
			break
		}

		if !ok && !isRegex {
			return
		}

		if pattern == child.path && !isRegex {
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

			for _, childNode := range n.children {
				queueTemp = append(queueTemp, childNode)
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
