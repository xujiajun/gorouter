package gorouter

import (
	"net/http"
	"context"
	"strings"
	"fmt"
	"regexp"
)

var (
	defaultPattern = `[\w]+`
	idPattern      = `[\d]+`
	idKey          = `id`
)

type (
	middlewareType func(next http.HandlerFunc) http.HandlerFunc

	Router struct {
		prefix     string
		middleware []middlewareType
		trees      map[string]*Tree
		notFound   http.HandlerFunc
	}
)

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

func (router *Router) PATCH(path string, handle http.HandlerFunc) {
	router.Handle(http.MethodPatch, path, handle)
}

func (router *Router) Group(prefix string) *Router {
	return &Router{
		prefix:     prefix,
		trees:      router.trees,
		middleware: router.middleware,
	}
}

func (router *Router) NotFoundFunc(handler http.HandlerFunc) {
	router.notFound = handler
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

	if router.prefix != "" {
		path = router.prefix + "/" + path
	}

	root.Add(path, handle, router.middleware...)
}

// GetParam return route param stored in r.
func GetParam(r *http.Request, key string) string {
	return GetAllParams(r)[key]
}

// contextKeyType is a private struct that is used for storing values in net.Context
type contextKeyType struct{}

// contextKey is the key that is used to store values in the net.Context for each request
var contextKey = contextKeyType{}

// paramsMapType is a private type that is used to store route params
type paramsMapType map[string]string

// GetAllParams return all route params stored in r.
func GetAllParams(r *http.Request) paramsMapType {
	values, ok := r.Context().Value(contextKey).(paramsMapType)
	if ok {
		return values
	}

	return paramsMapType{}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Path
	nodes := router.trees[r.Method].Find(requestUrl, 0)

	for _, node := range nodes {

		handler := node.handle
		path := node.path

		if handler != nil {
			if path == requestUrl {
				handle(w, r, handler, node.middleware)
				return
			}

			if path == requestUrl[1:] {
				handle(w, r, handler, node.middleware)
				return
			}
		}
	}

	if nodes == nil {
		res := strings.Split(requestUrl, "/")
		prefix := res[1]

		nodes := router.trees[r.Method].Find(prefix, 1)

		for _, node := range nodes {
			handler := node.handle

			if handler != nil && node.path != requestUrl {

				if matchParamsMap, ok := router.matchAndParse(requestUrl, node.path); ok {
					ctx := context.WithValue(r.Context(), contextKey, matchParamsMap)
					r = r.WithContext(ctx)
					handle(w, r, handler, node.middleware)
					return
				}
			}
		}
	}

	router.HandleNotFound(w, r, router.middleware)
}

func (router *Router) Use(middleware ...middlewareType) {
	if len(middleware) > 0 {
		router.middleware = append(router.middleware, middleware...)
	}
}

func (router *Router) HandleNotFound(w http.ResponseWriter, r *http.Request, middleware []middlewareType) {
	if router.notFound != nil {
		handle(w, r, router.notFound, middleware)
		return
	} else {
		http.NotFound(w, r)
	}
}

func handle(w http.ResponseWriter, r *http.Request, handler http.HandlerFunc, middleware []middlewareType) {
	var baseHandler = handler
	for _, m := range middleware {
		baseHandler = m(baseHandler)
	}
	baseHandler(w, r)
}

func (router *Router) Match(requestUrl string, path string) bool {
	_, ok := router.matchAndParse(requestUrl, path)
	return ok
}

func (router *Router) matchAndParse(requestUrl string, path string) (paramsMapType, bool) {
	res := strings.Split(path, "/")
	if res == nil {
		return nil, false
	}

	var (
		matchName []string
		sTemp     string
	)

	matchParams := make(paramsMapType)

	for _, str := range res {

		if str != "" {
			r := []byte(str)

			if string(r[0]) == "{" && string(r[len(r)-1]) == "}" {
				matchStr := string(r[1:len(r)-1])
				res := strings.Split(matchStr, ":")

				matchName = append(matchName, res[0])

				sTemp = sTemp + "/" + "(" + res[1] + ")"
			} else if string(r[0]) == ":" {
				matchStr := string(r)
				res := strings.Split(matchStr, ":")
				matchName = append(matchName, res[1])

				if res[1] == idKey {
					sTemp = sTemp + "/" + "(" + idPattern + ")"
				} else {
					sTemp = sTemp + "/" + "(" + defaultPattern + ")"
				}
			} else {
				sTemp = sTemp + "/" + str
			}
		}
	}

	pattern := sTemp

	re := regexp.MustCompile(pattern)
	submatch := re.FindSubmatch([]byte(requestUrl))

	if submatch != nil {
		if string(submatch[0]) == requestUrl {
			submatch = submatch[1:]
			for k, v := range submatch {
				matchParams[matchName[k]] = string(v)
			}
			return matchParams, true
		}
	}

	return nil, false
}
