package gorouter_test

import (
	"github.com/julienschmidt/httprouter"
	"github.com/xujiajun/gorouter"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	//triemux "github.com/teambition/trie-mux/mux"
	"github.com/gorilla/mux"
)

type route struct {
	method string
	path   string
}

// http://developer.github.com/v3/
var githubAPI = []route{
	// OAuth Authorizations
	{"GET", "/authorizations"},
	{"GET", "/authorizations/:id"},
	{"POST", "/authorizations"},
	//{"PUT", "/authorizations/clients/:client_id"},
	//{"PATCH", "/authorizations/:id"},
	{"DELETE", "/authorizations/:id"},
	{"GET", "/applications/:client_id/tokens/:access_token"},
	{"DELETE", "/applications/:client_id/tokens"},
	{"DELETE", "/applications/:client_id/tokens/:access_token"},

	// Activity
	{"GET", "/events"},
	{"GET", "/repos/:owner/:repo/events"},
	{"GET", "/networks/:owner/:repo/events"},
	{"GET", "/orgs/:org/events"},
	{"GET", "/users/:user/received_events"},
	{"GET", "/users/:user/received_events/public"},
	{"GET", "/users/:user/events"},
	{"GET", "/users/:user/events/public"},
	{"GET", "/users/:user/events/orgs/:org"},
	{"GET", "/feeds"},
	{"GET", "/notifications"},
	{"GET", "/repos/:owner/:repo/notifications"},
	{"PUT", "/notifications"},
	{"PUT", "/repos/:owner/:repo/notifications"},
	{"GET", "/notifications/threads/:id"},
	//{"PATCH", "/notifications/threads/:id"},
	{"GET", "/notifications/threads/:id/subscription"},
	{"PUT", "/notifications/threads/:id/subscription"},
	{"DELETE", "/notifications/threads/:id/subscription"},
	{"GET", "/repos/:owner/:repo/stargazers"},
	{"GET", "/users/:user/starred"},
	{"GET", "/user/starred"},
	{"GET", "/user/starred/:owner/:repo"},
	{"PUT", "/user/starred/:owner/:repo"},
	{"DELETE", "/user/starred/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/subscribers"},
	{"GET", "/users/:user/subscriptions"},
	{"GET", "/user/subscriptions"},
	{"GET", "/repos/:owner/:repo/subscription"},
	{"PUT", "/repos/:owner/:repo/subscription"},
	{"DELETE", "/repos/:owner/:repo/subscription"},
	{"GET", "/user/subscriptions/:owner/:repo"},
	{"PUT", "/user/subscriptions/:owner/:repo"},
	{"DELETE", "/user/subscriptions/:owner/:repo"},

	// Gists
	{"GET", "/users/:user/gists"},
	{"GET", "/gists"},
	//{"GET", "/gists/public"},
	//{"GET", "/gists/starred"},
	{"GET", "/gists/:id"},
	{"POST", "/gists"},
	//{"PATCH", "/gists/:id"},
	{"PUT", "/gists/:id/star"},
	{"DELETE", "/gists/:id/star"},
	{"GET", "/gists/:id/star"},
	{"POST", "/gists/:id/forks"},
	{"DELETE", "/gists/:id"},

	// Git Data
	{"GET", "/repos/:owner/:repo/git/blobs/:sha"},
	{"POST", "/repos/:owner/:repo/git/blobs"},
	{"GET", "/repos/:owner/:repo/git/commits/:sha"},
	{"POST", "/repos/:owner/:repo/git/commits"},
	//{"GET", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/refs"},
	{"POST", "/repos/:owner/:repo/git/refs"},
	//{"PATCH", "/repos/:owner/:repo/git/refs/*ref"},
	//{"DELETE", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/tags/:sha"},
	{"POST", "/repos/:owner/:repo/git/tags"},
	{"GET", "/repos/:owner/:repo/git/trees/:sha"},
	{"POST", "/repos/:owner/:repo/git/trees"},

	// Issues
	{"GET", "/issues"},
	{"GET", "/user/issues"},
	{"GET", "/orgs/:org/issues"},
	{"GET", "/repos/:owner/:repo/issues"},
	{"GET", "/repos/:owner/:repo/issues/:number"},
	{"POST", "/repos/:owner/:repo/issues"},
	//{"PATCH", "/repos/:owner/:repo/issues/:number"},
	{"GET", "/repos/:owner/:repo/assignees"},
	{"GET", "/repos/:owner/:repo/assignees/:assignee"},
	{"GET", "/repos/:owner/:repo/issues/:number/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments/:id"},
	{"POST", "/repos/:owner/:repo/issues/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/issues/comments/:id"},
	//{"DELETE", "/repos/:owner/:repo/issues/comments/:id"},
	{"GET", "/repos/:owner/:repo/issues/:number/events"},
	//{"GET", "/repos/:owner/:repo/issues/events"},
	//{"GET", "/repos/:owner/:repo/issues/events/:id"},
	{"GET", "/repos/:owner/:repo/labels"},
	{"GET", "/repos/:owner/:repo/labels/:name"},
	{"POST", "/repos/:owner/:repo/labels"},
	//{"PATCH", "/repos/:owner/:repo/labels/:name"},
	{"DELETE", "/repos/:owner/:repo/labels/:name"},
	{"GET", "/repos/:owner/:repo/issues/:number/labels"},
	{"POST", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name"},
	{"PUT", "/repos/:owner/:repo/issues/:number/labels"},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones/:number/labels"},
	{"GET", "/repos/:owner/:repo/milestones"},
	{"GET", "/repos/:owner/:repo/milestones/:number"},
	{"POST", "/repos/:owner/:repo/milestones"},
	//{"PATCH", "/repos/:owner/:repo/milestones/:number"},
	{"DELETE", "/repos/:owner/:repo/milestones/:number"},

	// Miscellaneous
	{"GET", "/emojis"},
	{"GET", "/gitignore/templates"},
	{"GET", "/gitignore/templates/:name"},
	{"POST", "/markdown"},
	{"POST", "/markdown/raw"},
	{"GET", "/meta"},
	{"GET", "/rate_limit"},

	// Organizations
	{"GET", "/users/:user/orgs"},
	{"GET", "/user/orgs"},
	{"GET", "/orgs/:org"},
	//{"PATCH", "/orgs/:org"},
	{"GET", "/orgs/:org/members"},
	{"GET", "/orgs/:org/members/:user"},
	{"DELETE", "/orgs/:org/members/:user"},
	{"GET", "/orgs/:org/public_members"},
	{"GET", "/orgs/:org/public_members/:user"},
	{"PUT", "/orgs/:org/public_members/:user"},
	{"DELETE", "/orgs/:org/public_members/:user"},
	{"GET", "/orgs/:org/teams"},
	{"GET", "/teams/:id"},
	{"POST", "/orgs/:org/teams"},
	//{"PATCH", "/teams/:id"},
	{"DELETE", "/teams/:id"},
	{"GET", "/teams/:id/members"},
	{"GET", "/teams/:id/members/:user"},
	{"PUT", "/teams/:id/members/:user"},
	{"DELETE", "/teams/:id/members/:user"},
	{"GET", "/teams/:id/repos"},
	{"GET", "/teams/:id/repos/:owner/:repo"},
	{"PUT", "/teams/:id/repos/:owner/:repo"},
	{"DELETE", "/teams/:id/repos/:owner/:repo"},
	{"GET", "/user/teams"},

	// Pull Requests
	{"GET", "/repos/:owner/:repo/pulls"},
	{"GET", "/repos/:owner/:repo/pulls/:number"},
	{"POST", "/repos/:owner/:repo/pulls"},
	//{"PATCH", "/repos/:owner/:repo/pulls/:number"},
	{"GET", "/repos/:owner/:repo/pulls/:number/commits"},
	{"GET", "/repos/:owner/:repo/pulls/:number/files"},
	{"GET", "/repos/:owner/:repo/pulls/:number/merge"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/merge"},
	{"GET", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments/:number"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/comments"},
	//{"PATCH", "/repos/:owner/:repo/pulls/comments/:number"},
	//{"DELETE", "/repos/:owner/:repo/pulls/comments/:number"},

	// Repositories
	{"GET", "/user/repos"},
	{"GET", "/users/:user/repos"},
	{"GET", "/orgs/:org/repos"},
	{"GET", "/repositories"},
	{"POST", "/user/repos"},
	{"POST", "/orgs/:org/repos"},
	{"GET", "/repos/:owner/:repo"},
	//{"PATCH", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/contributors"},
	{"GET", "/repos/:owner/:repo/languages"},
	{"GET", "/repos/:owner/:repo/teams"},
	{"GET", "/repos/:owner/:repo/tags"},
	{"GET", "/repos/:owner/:repo/branches"},
	{"GET", "/repos/:owner/:repo/branches/:branch"},
	{"DELETE", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/collaborators"},
	{"GET", "/repos/:owner/:repo/collaborators/:user"},
	{"PUT", "/repos/:owner/:repo/collaborators/:user"},
	{"DELETE", "/repos/:owner/:repo/collaborators/:user"},
	{"GET", "/repos/:owner/:repo/comments"},
	{"GET", "/repos/:owner/:repo/commits/:sha/comments"},
	{"POST", "/repos/:owner/:repo/commits/:sha/comments"},
	{"GET", "/repos/:owner/:repo/comments/:id"},
	//{"PATCH", "/repos/:owner/:repo/comments/:id"},
	{"DELETE", "/repos/:owner/:repo/comments/:id"},
	{"GET", "/repos/:owner/:repo/commits"},
	{"GET", "/repos/:owner/:repo/commits/:sha"},
	{"GET", "/repos/:owner/:repo/readme"},
	//{"GET", "/repos/:owner/:repo/contents/*path"},
	//{"PUT", "/repos/:owner/:repo/contents/*path"},
	//{"DELETE", "/repos/:owner/:repo/contents/*path"},
	//{"GET", "/repos/:owner/:repo/:archive_format/:ref"},
	{"GET", "/repos/:owner/:repo/keys"},
	{"GET", "/repos/:owner/:repo/keys/:id"},
	{"POST", "/repos/:owner/:repo/keys"},
	//{"PATCH", "/repos/:owner/:repo/keys/:id"},
	{"DELETE", "/repos/:owner/:repo/keys/:id"},
	{"GET", "/repos/:owner/:repo/downloads"},
	{"GET", "/repos/:owner/:repo/downloads/:id"},
	{"DELETE", "/repos/:owner/:repo/downloads/:id"},
	{"GET", "/repos/:owner/:repo/forks"},
	{"POST", "/repos/:owner/:repo/forks"},
	{"GET", "/repos/:owner/:repo/hooks"},
	{"GET", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks"},
	//{"PATCH", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks/:id/tests"},
	{"DELETE", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/merges"},
	{"GET", "/repos/:owner/:repo/releases"},
	{"GET", "/repos/:owner/:repo/releases/:id"},
	{"POST", "/repos/:owner/:repo/releases"},
	//{"PATCH", "/repos/:owner/:repo/releases/:id"},
	{"DELETE", "/repos/:owner/:repo/releases/:id"},
	{"GET", "/repos/:owner/:repo/releases/:id/assets"},
	{"GET", "/repos/:owner/:repo/stats/contributors"},
	{"GET", "/repos/:owner/:repo/stats/commit_activity"},
	{"GET", "/repos/:owner/:repo/stats/code_frequency"},
	{"GET", "/repos/:owner/:repo/stats/participation"},
	{"GET", "/repos/:owner/:repo/stats/punch_card"},
	{"GET", "/repos/:owner/:repo/statuses/:ref"},
	{"POST", "/repos/:owner/:repo/statuses/:ref"},

	// Search
	{"GET", "/search/repositories"},
	{"GET", "/search/code"},
	{"GET", "/search/issues"},
	{"GET", "/search/users"},
	{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword"},
	{"GET", "/legacy/repos/search/:keyword"},
	{"GET", "/legacy/user/search/:keyword"},
	{"GET", "/legacy/user/email/:email"},

	// Users
	{"GET", "/users/:user"},
	{"GET", "/user"},
	//{"PATCH", "/user"},
	{"GET", "/users"},
	{"GET", "/user/emails"},
	{"POST", "/user/emails"},
	{"DELETE", "/user/emails"},
	{"GET", "/users/:user/followers"},
	{"GET", "/user/followers"},
	{"GET", "/users/:user/following"},
	{"GET", "/user/following"},
	{"GET", "/user/following/:user"},
	{"GET", "/users/:user/following/:target_user"},
	{"PUT", "/user/following/:user"},
	{"DELETE", "/user/following/:user"},
	{"GET", "/users/:user/keys"},
	{"GET", "/user/keys"},
	{"GET", "/user/keys/:id"},
	{"POST", "/user/keys"},
	//{"PATCH", "/user/keys/:id"},
	{"DELETE", "/user/keys/:id"},
}

var githubAPI2 = []route{
	// OAuth Authorizations
	{"GET", "/authorizations"},
	{"GET", "/authorizations/{id:[0-9]+}"},
	{"POST", "/authorizations"},
	//{"PUT", "/authorizations/clients/{client_id:[0-9]+}"},
	//{"PATCH", "/authorizations/{id:[0-9]+}"},
	{"DELETE", "/authorizations/{id:[0-9]+}"},
	{"GET", "/applications/{client_id:[0-9]+}/tokens/{access_token}"},
	{"DELETE", "/applications/{client_id:[0-9]+}/tokens"},
	{"DELETE", "/applications/{client_id:[0-9]+}/tokens/{access_token}"},

	// Activity
	{"GET", "/events"},
	{"GET", "/repos/{owner}/{repo}/events"},
	{"GET", "/networks/{owner}/{repo}/events"},
	{"GET", "/orgs/{org}/events"},
	{"GET", "/users/{user}/received_events"},
	{"GET", "/users/{user}/received_events/public"},
	{"GET", "/users/{user}/events"},
	{"GET", "/users/{user}/events/public"},
	{"GET", "/users/{user}/events/orgs/{org}"},
	{"GET", "/feeds"},
	{"GET", "/notifications"},
	{"GET", "/repos/{owner}/{repo}/notifications"},
	{"PUT", "/notifications"},
	{"PUT", "/repos/{owner}/{repo}/notifications"},
	{"GET", "/notifications/threads/{id:[0-9]+}"},
	//{"PATCH", "/notifications/threads/{id:[0-9]+}"},
	{"GET", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"PUT", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"DELETE", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"GET", "/repos/{owner}/{repo}/stargazers"},
	{"GET", "/users/{user}/starred"},
	{"GET", "/user/starred"},
	{"GET", "/user/starred/{owner}/{repo}"},
	{"PUT", "/user/starred/{owner}/{repo}"},
	{"DELETE", "/user/starred/{owner}/{repo}"},
	{"GET", "/repos/{owner}/{repo}/subscribers"},
	{"GET", "/users/{user}/subscriptions"},
	{"GET", "/user/subscriptions"},
	{"GET", "/repos/{owner}/{repo}/subscription"},
	{"PUT", "/repos/{owner}/{repo}/subscription"},
	{"DELETE", "/repos/{owner}/{repo}/subscription"},
	{"GET", "/user/subscriptions/{owner}/{repo}"},
	{"PUT", "/user/subscriptions/{owner}/{repo}"},
	{"DELETE", "/user/subscriptions/{owner}/{repo}"},

	// Gists
	{"GET", "/users/{user}/gists"},
	{"GET", "/gists"},
	//{"GET", "/gists/public"},
	//{"GET", "/gists/starred"},
	{"GET", "/gists/{id:[0-9]+}"},
	{"POST", "/gists"},
	//{"PATCH", "/gists/{id:[0-9]+}"},
	{"PUT", "/gists/{id:[0-9]+}/star"},
	{"DELETE", "/gists/{id:[0-9]+}/star"},
	{"GET", "/gists/{id:[0-9]+}/star"},
	{"POST", "/gists/{id:[0-9]+}/forks"},
	{"DELETE", "/gists/{id:[0-9]+}"},

	// Git Data
	{"GET", "/repos/{owner}/{repo}/git/blobs/{sha}"},
	{"POST", "/repos/{owner}/{repo}/git/blobs"},
	{"GET", "/repos/{owner}/{repo}/git/commits/{sha}"},
	{"POST", "/repos/{owner}/{repo}/git/commits"},
	//{"GET", "/repos/{owner}/{repo}/git/refs/*ref"},
	{"GET", "/repos/{owner}/{repo}/git/refs"},
	{"POST", "/repos/{owner}/{repo}/git/refs"},
	//{"PATCH", "/repos/{owner}/{repo}/git/refs/*ref"},
	//{"DELETE", "/repos/{owner}/{repo}/git/refs/*ref"},
	{"GET", "/repos/{owner}/{repo}/git/tags/{sha}"},
	{"POST", "/repos/{owner}/{repo}/git/tags"},
	{"GET", "/repos/{owner}/{repo}/git/trees/{sha}"},
	{"POST", "/repos/{owner}/{repo}/git/trees"},

	// Issues
	{"GET", "/issues"},
	{"GET", "/user/issues"},
	{"GET", "/orgs/{org}/issues"},
	{"GET", "/repos/{owner}/{repo}/issues"},
	{"GET", "/repos/{owner}/{repo}/issues/{number}"},
	{"POST", "/repos/{owner}/{repo}/issues"},
	//{"PATCH", "/repos/{owner}/{repo}/issues/{number}"},
	{"GET", "/repos/{owner}/{repo}/assignees"},
	{"GET", "/repos/{owner}/{repo}/assignees/{assignee}"},
	{"GET", "/repos/{owner}/{repo}/issues/{number}/comments"},
	//{"GET", "/repos/{owner}/{repo}/issues/comments"},
	//{"GET", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/issues/{number}/comments"},
	//{"PATCH", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	//{"DELETE", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/issues/{number}/events"},
	//{"GET", "/repos/{owner}/{repo}/issues/events"},
	//{"GET", "/repos/{owner}/{repo}/issues/events/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/labels"},
	{"GET", "/repos/{owner}/{repo}/labels/{name}"},
	{"POST", "/repos/{owner}/{repo}/labels"},
	//{"PATCH", "/repos/{owner}/{repo}/labels/{name}"},
	{"DELETE", "/repos/{owner}/{repo}/labels/{name}"},
	{"GET", "/repos/{owner}/{repo}/issues/{number}/labels"},
	{"POST", "/repos/{owner}/{repo}/issues/{number}/labels"},
	{"DELETE", "/repos/{owner}/{repo}/issues/{number}/labels/{name}"},
	{"PUT", "/repos/{owner}/{repo}/issues/{number}/labels"},
	{"DELETE", "/repos/{owner}/{repo}/issues/{number}/labels"},
	{"GET", "/repos/{owner}/{repo}/milestones/{number}/labels"},
	{"GET", "/repos/{owner}/{repo}/milestones"},
	{"GET", "/repos/{owner}/{repo}/milestones/{number}"},
	{"POST", "/repos/{owner}/{repo}/milestones"},
	//{"PATCH", "/repos/{owner}/{repo}/milestones/{number}"},
	{"DELETE", "/repos/{owner}/{repo}/milestones/{number}"},

	// Miscellaneous
	{"GET", "/emojis"},
	{"GET", "/gitignore/templates"},
	{"GET", "/gitignore/templates/{name}"},
	{"POST", "/markdown"},
	{"POST", "/markdown/raw"},
	{"GET", "/meta"},
	{"GET", "/rate_limit"},

	// Organizations
	{"GET", "/users/{user}/orgs"},
	{"GET", "/user/orgs"},
	{"GET", "/orgs/{org}"},
	//{"PATCH", "/orgs/{org}"},
	{"GET", "/orgs/{org}/members"},
	{"GET", "/orgs/{org}/members/{user}"},
	{"DELETE", "/orgs/{org}/members/{user}"},
	{"GET", "/orgs/{org}/public_members"},
	{"GET", "/orgs/{org}/public_members/{user}"},
	{"PUT", "/orgs/{org}/public_members/{user}"},
	{"DELETE", "/orgs/{org}/public_members/{user}"},
	{"GET", "/orgs/{org}/teams"},
	{"GET", "/teams/{id:[0-9]+}"},
	{"POST", "/orgs/{org}/teams"},
	//{"PATCH", "/teams/{id:[0-9]+}"},
	{"DELETE", "/teams/{id:[0-9]+}"},
	{"GET", "/teams/{id:[0-9]+}/members"},
	{"GET", "/teams/{id:[0-9]+}/members/{user}"},
	{"PUT", "/teams/{id:[0-9]+}/members/{user}"},
	{"DELETE", "/teams/{id:[0-9]+}/members/{user}"},
	{"GET", "/teams/{id:[0-9]+}/repos"},
	{"GET", "/teams/{id:[0-9]+}/repos/{owner}/{repo}"},
	{"PUT", "/teams/{id:[0-9]+}/repos/{owner}/{repo}"},
	{"DELETE", "/teams/{id:[0-9]+}/repos/{owner}/{repo}"},
	{"GET", "/user/teams"},

	// Pull Requests
	{"GET", "/repos/{owner}/{repo}/pulls"},
	{"GET", "/repos/{owner}/{repo}/pulls/{number}"},
	{"POST", "/repos/{owner}/{repo}/pulls"},
	//{"PATCH", "/repos/{owner}/{repo}/pulls/{number}"},
	{"GET", "/repos/{owner}/{repo}/pulls/{number}/commits"},
	{"GET", "/repos/{owner}/{repo}/pulls/{number}/files"},
	{"GET", "/repos/{owner}/{repo}/pulls/{number}/merge"},
	{"PUT", "/repos/{owner}/{repo}/pulls/{number}/merge"},
	{"GET", "/repos/{owner}/{repo}/pulls/{number}/comments"},
	//{"GET", "/repos/{owner}/{repo}/pulls/comments"},
	//{"GET", "/repos/{owner}/{repo}/pulls/comments/{number}"},
	{"PUT", "/repos/{owner}/{repo}/pulls/{number}/comments"},
	//{"PATCH", "/repos/{owner}/{repo}/pulls/comments/{number}"},
	//{"DELETE", "/repos/{owner}/{repo}/pulls/comments/{number}"},

	// Repositories
	{"GET", "/user/repos"},
	{"GET", "/users/{user}/repos"},
	{"GET", "/orgs/{org}/repos"},
	{"GET", "/repositories"},
	{"POST", "/user/repos"},
	{"POST", "/orgs/{org}/repos"},
	{"GET", "/repos/{owner}/{repo}"},
	//{"PATCH", "/repos/{owner}/{repo}"},
	{"GET", "/repos/{owner}/{repo}/contributors"},
	{"GET", "/repos/{owner}/{repo}/languages"},
	{"GET", "/repos/{owner}/{repo}/teams"},
	{"GET", "/repos/{owner}/{repo}/tags"},
	{"GET", "/repos/{owner}/{repo}/branches"},
	{"GET", "/repos/{owner}/{repo}/branches/{branch}"},
	{"DELETE", "/repos/{owner}/{repo}"},
	{"GET", "/repos/{owner}/{repo}/collaborators"},
	{"GET", "/repos/{owner}/{repo}/collaborators/{user}"},
	{"PUT", "/repos/{owner}/{repo}/collaborators/{user}"},
	{"DELETE", "/repos/{owner}/{repo}/collaborators/{user}"},
	{"GET", "/repos/{owner}/{repo}/comments"},
	{"GET", "/repos/{owner}/{repo}/commits/{sha}/comments"},
	{"POST", "/repos/{owner}/{repo}/commits/{sha}/comments"},
	{"GET", "/repos/{owner}/{repo}/comments/{id:[0-9]+}"},
	//{"PATCH", "/repos/{owner}/{repo}/comments/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner}/{repo}/comments/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/commits"},
	{"GET", "/repos/{owner}/{repo}/commits/{sha}"},
	{"GET", "/repos/{owner}/{repo}/readme"},
	//{"GET", "/repos/{owner}/{repo}/contents/*path"},
	//{"PUT", "/repos/{owner}/{repo}/contents/*path"},
	//{"DELETE", "/repos/{owner}/{repo}/contents/*path"},
	//{"GET", "/repos/{owner}/{repo}/:archive_format/:ref"},
	{"GET", "/repos/{owner}/{repo}/keys"},
	{"GET", "/repos/{owner}/{repo}/keys/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/keys"},
	//{"PATCH", "/repos/{owner}/{repo}/keys/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner}/{repo}/keys/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/downloads"},
	{"GET", "/repos/{owner}/{repo}/downloads/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner}/{repo}/downloads/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/forks"},
	{"POST", "/repos/{owner}/{repo}/forks"},
	{"GET", "/repos/{owner}/{repo}/hooks"},
	{"GET", "/repos/{owner}/{repo}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/hooks"},
	//{"PATCH", "/repos/{owner}/{repo}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/hooks/{id:[0-9]+}/tests"},
	{"DELETE", "/repos/{owner}/{repo}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/merges"},
	{"GET", "/repos/{owner}/{repo}/releases"},
	{"GET", "/repos/{owner}/{repo}/releases/{id:[0-9]+}"},
	{"POST", "/repos/{owner}/{repo}/releases"},
	//{"PATCH", "/repos/{owner}/{repo}/releases/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner}/{repo}/releases/{id:[0-9]+}"},
	{"GET", "/repos/{owner}/{repo}/releases/{id:[0-9]+}/assets"},
	{"GET", "/repos/{owner}/{repo}/stats/contributors"},
	{"GET", "/repos/{owner}/{repo}/stats/commit_activity"},
	{"GET", "/repos/{owner}/{repo}/stats/code_frequency"},
	{"GET", "/repos/{owner}/{repo}/stats/participation"},
	{"GET", "/repos/{owner}/{repo}/stats/punch_card"},
	{"GET", "/repos/{owner}/{repo}/statuses/{ref}"},
	{"POST", "/repos/{owner}/{repo}/statuses/{ref}"},

	// Search
	{"GET", "/search/repositories"},
	{"GET", "/search/code"},
	{"GET", "/search/issues"},
	{"GET", "/search/users"},
	{"GET", "/legacy/issues/search/{owner}/{repo}sitory/{state}/{keyword}"},
	{"GET", "/legacy/repos/search/{keyword}"},
	{"GET", "/legacy/user/search/{keyword}"},
	{"GET", "/legacy/user/email/{email}"},

	// Users
	{"GET", "/users/{user}"},
	{"GET", "/user"},
	//{"PATCH", "/user"},
	{"GET", "/users"},
	{"GET", "/user/emails"},
	{"POST", "/user/emails"},
	{"DELETE", "/user/emails"},
	{"GET", "/users/{user}/followers"},
	{"GET", "/user/followers"},
	{"GET", "/users/{user}/following"},
	{"GET", "/user/following"},
	{"GET", "/user/following/{user}"},
	{"GET", "/users/{user}/following/{target_user}"},
	{"PUT", "/user/following/{user}"},
	{"DELETE", "/user/following/{user}"},
	{"GET", "/users/{user}/keys"},
	{"GET", "/user/keys"},
	{"GET", "/user/keys/{id:[0-9]+}"},
	{"POST", "/user/keys"},
	//{"PATCH", "/user/keys/{id:[0-9]+}"},
	{"DELETE", "/user/keys/{id:[0-9]+}"},
}

func calcMem(name string, load func()) {
	m := new(runtime.MemStats)

	// before
	runtime.GC()
	runtime.ReadMemStats(m)
	before := m.HeapAlloc

	load()

	// after
	runtime.GC()
	runtime.ReadMemStats(m)
	after := m.HeapAlloc
	println("   "+name+":", after-before, "Bytes")
}

var (
	httpRouter http.Handler
	goRouter   http.Handler
	muxRouter  http.Handler
	//trieMuxRouter http.Handler
)

func init() {
	println("GithubAPI Routes:", len(githubAPI))
	println("GithubAPI2 Routes:", len(githubAPI2))
	calcMem("HttpRouter", func() {
		router := httprouter.New()
		handler := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.method, route.path, handler)
		}
		httpRouter = router
	})

	calcMem("GoRouter", func() {
		router := gorouter.New()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.method, route.path, handler)
		}
		goRouter = router
	})

	// go version >= 1.8
	//calcMem("trie-mux", func() {
	//	router := triemux.New()
	//	handler := func(w http.ResponseWriter, _ *http.Request, _ triemux.Params) {
	//		w.WriteHeader(204)
	//	}
	//	for _, route := range githubAPI {
	//		router.Handle(route.method, route.path, handler)
	//	}
	//	trieMuxRouter = router
	//})

	calcMem("MuxRouter", func() {
		router := mux.NewRouter()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI2 {
			router.HandleFunc(route.path, handler).Methods(route.method)
		}
		muxRouter = router
	})
}

func benchRoutes(b *testing.B, router http.Handler, routes []route) {
	b.N = 10000
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, route := range routes {
			router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(route.method, route.path, nil))
		}
	}
}

// go version >= 1.8
//func BenchmarkTrieMuxRouter(b *testing.B) {
//	benchRoutes(b, trieMuxRouter, githubAPI)
//}

func BenchmarkHttpRouter(b *testing.B) {
	benchRoutes(b, httpRouter, githubAPI)
}

func BenchmarkGoRouter(b *testing.B) {
	benchRoutes(b, goRouter, githubAPI)
}

func BenchmarkMuxRouter(b *testing.B) {
	benchRoutes(b, muxRouter, githubAPI2)
}
