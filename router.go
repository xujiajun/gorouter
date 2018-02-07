package gorouter

import (
	"net/http"
	"context"
	"strings"
	"fmt"
)

type Router struct {
	trees map[string]*Tree
}

func New() *Router {
	return &Router{
		trees: make(map[string]*Tree),
	}
}

func (router *Router) GET(path string, handle http.HandlerFunc) {
	router.Handle(http.MethodGet, path, handle)
}

func (router *Router) POST(path string, handle http.HandlerFunc) {
	router.Handle(http.MethodPost, path, handle)
}

func (router *Router) DELETE(path string, handle http.HandlerFunc) {
	router.Handle(http.MethodDelete, path, handle)
}

func (router *Router) PUT(path string, handle http.HandlerFunc) {
	router.Handle(http.MethodPut, path, handle)
}

func (router *Router) Handle(method string, path string, handle http.HandlerFunc) {
	if method == "" {
		panic(fmt.Errorf("invalid method"))
	}

	if router.trees == nil {
		router.trees = make(map[string]*Tree)
	}

	root := router.trees[method]
	if root == nil {
		root = NewTree()
		router.trees[method] = root
	}

	root.Add(path, handle)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestUrl := r.URL.Path
	nodes := router.trees[r.Method].Find(requestUrl)

	for _, node := range nodes {

		handler := node.Handle

		if handler != nil {
			if node.Path == requestUrl {
				handler(w, r)
				return
			}
		}
	}

	if nodes == nil {
		// find again
		res := strings.Split(requestUrl, "/")
		prefix := "/" + res[0]

		nodes := router.trees[r.Method].Find(prefix)

		for _, node := range nodes {
			handler := node.Handle

			if handler != nil && node.Path != requestUrl {
				isMatch, matchParams := Match(requestUrl, node.Path)
				if isMatch {
					for k, v := range matchParams {
						ctx := context.WithValue(r.Context(), k, v)
						r = r.WithContext(ctx)
					}

					handler(w, r)
					return
				}
			}
		}

		http.NotFound(w, r)
		return
	}

	http.NotFound(w, r)
	return
}
