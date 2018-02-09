package gorouter

import (
	"net/http"
	"strings"
)

type (
	Tree struct {
		root *Node
	}

	Node struct {
		key        string
		path       string
		handle     http.HandlerFunc
		depth      int
		children   map[string]*Node
		isLeaf     bool
		middleware []middlewareType
	}
)

func NewNode(key string, depth int) *Node {
	return &Node{
		key:      key,
		depth:    depth,
		children: make(map[string]*Node),
	}
}

func NewTree() *Tree {
	return &Tree{
		root: NewNode("/", 1),
	}
}

func (tree *Tree) Add(pattern string, handle http.HandlerFunc, middleware ...middlewareType) {
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
	parent.isLeaf = true
	parent.path = pattern

}

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
