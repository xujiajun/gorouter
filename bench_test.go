package gorouter_test

import (
	beegomux "github.com/beego/mux"
	"github.com/go-chi/chi"
	"github.com/go-zoo/bone"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	triemux "github.com/teambition/trie-mux/mux"
	"github.com/xujiajun/gorouter"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
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
	{"GET", "/applications/{client_id:[0-9]+}/tokens/{access_token:\\w+}"},
	{"DELETE", "/applications/{client_id:[0-9]+}/tokens"},
	{"DELETE", "/applications/{client_id:[0-9]+}/tokens/{access_token:\\w+}"},

	// Activity
	{"GET", "/events"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/events"},
	{"GET", "/networks/{owner:\\w+}/{repo:\\w+}/events"},
	{"GET", "/orgs/{org:\\w+}/events"},
	{"GET", "/users/{user:\\w+}/received_events"},
	{"GET", "/users/{user:\\w+}/received_events/public"},
	{"GET", "/users/{user:\\w+}/events"},
	{"GET", "/users/{user:\\w+}/events/public"},
	{"GET", "/users/{user:\\w+}/events/orgs/{org:\\w+}"},
	{"GET", "/feeds"},
	{"GET", "/notifications"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/notifications"},
	{"PUT", "/notifications"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/notifications"},
	{"GET", "/notifications/threads/{id:[0-9]+}"},
	//{"PATCH", "/notifications/threads/{id:[0-9]+}"},
	{"GET", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"PUT", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"DELETE", "/notifications/threads/{id:[0-9]+}/subscription"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stargazers"},
	{"GET", "/users/{user:\\w+}/starred"},
	{"GET", "/user/starred"},
	{"GET", "/user/starred/{owner:\\w+}/{repo:\\w+}"},
	{"PUT", "/user/starred/{owner:\\w+}/{repo:\\w+}"},
	{"DELETE", "/user/starred/{owner:\\w+}/{repo:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/subscribers"},
	{"GET", "/users/{user:\\w+}/subscriptions"},
	{"GET", "/user/subscriptions"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/subscription"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/subscription"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/subscription"},
	{"GET", "/user/subscriptions/{owner:\\w+}/{repo:\\w+}"},
	{"PUT", "/user/subscriptions/{owner:\\w+}/{repo:\\w+}"},
	{"DELETE", "/user/subscriptions/{owner:\\w+}/{repo:\\w+}"},

	// Gists
	{"GET", "/users/{user:\\w+}/gists"},
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
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/git/blobs/{sha:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/git/blobs"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/git/commits/{sha:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/git/commits"},
	//{"GET", "/repos/{owner}/{repo}/git/refs/*ref"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/git/refs"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/git/refs"},
	//{"PATCH", "/repos/{owner}/{repo}/git/refs/*ref"},
	//{"DELETE", "/repos/{owner}/{repo}/git/refs/*ref"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/git/tags/{sha:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/git/tags"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/git/trees/{sha:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/git/trees"},

	// Issues
	{"GET", "/issues"},
	{"GET", "/user/issues"},
	{"GET", "/orgs/{org:\\w+}/issues"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/issues"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{number:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/issues"},
	//{"PATCH", "/repos/{owner}/{repo}/issues/{number}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/assignees"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/assignees/{assignee:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{number:\\w+}/comments"},
	//{"GET", "/repos/{owner}/{repo}/issues/comments"},
	//{"GET", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{number:\\w+}/comments"},
	//{"PATCH", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	//{"DELETE", "/repos/{owner}/{repo}/issues/comments/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{number:\\w+}/events"},
	//{"GET", "/repos/{owner}/{repo}/issues/events"},
	//{"GET", "/repos/{owner}/{repo}/issues/events/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/labels"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/labels/{name:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/labels"},
	//{"PATCH", "/repos/{owner}/{repo}/labels/{name}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/labels/{name:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{name:\\w+}/labels"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{name:\\w+}/labels"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{name:\\w+}/labels/{name:\\w+}"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{name:\\w+}/labels"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/issues/{name:\\w+}/labels"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/milestones/{name:\\w+}/labels"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/milestones"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/milestones/{name:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/milestones"},
	//{"PATCH", "/repos/{owner}/{repo}/milestones/{number}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/milestones/{number:\\w+}"},

	// Miscellaneous
	{"GET", "/emojis"},
	{"GET", "/gitignore/templates"},
	{"GET", "/gitignore/templates/{name:\\w+}"},
	{"POST", "/markdown"},
	{"POST", "/markdown/raw"},
	{"GET", "/meta"},
	{"GET", "/rate_limit"},

	// Organizations
	{"GET", "/users/{user:\\w+}/orgs"},
	{"GET", "/user/orgs"},
	{"GET", "/orgs/{org:\\w+}"},
	//{"PATCH", "/orgs/{org}"},
	{"GET", "/orgs/{org:\\w+}/members"},
	{"GET", "/orgs/{org:\\w+}/members/{user:\\w+}"},
	{"DELETE", "/orgs/{org:\\w+}/members/{user:\\w+}"},
	{"GET", "/orgs/{org:\\w+}/public_members"},
	{"GET", "/orgs/{org:\\w+}/public_members/{user:\\w+}"},
	{"PUT", "/orgs/{org:\\w+}/public_members/{user:\\w+}"},
	{"DELETE", "/orgs/{org:\\w+}/public_members/{user:\\w+}"},
	{"GET", "/orgs/{org:\\w+}/teams"},
	{"GET", "/teams/{id:[0-9]+}"},
	{"POST", "/orgs/{org:\\w+}/teams"},
	//{"PATCH", "/teams/{id:[0-9]+}"},
	{"DELETE", "/teams/{id:[0-9]+}"},
	{"GET", "/teams/{id:[0-9]+}/members"},
	{"GET", "/teams/{id:[0-9]+}/members/{user:\\w+}"},
	{"PUT", "/teams/{id:[0-9]+}/members/{user:\\w+}"},
	{"DELETE", "/teams/{id:[0-9]+}/members/{user:\\w+}"},
	{"GET", "/teams/{id:[0-9]+}/repos"},
	{"GET", "/teams/{id:[0-9]+}/repos/{owner:\\w+}/{repo:\\w+}"},
	{"PUT", "/teams/{id:[0-9]+}/repos/{owner:\\w+}/{repo:\\w+}"},
	{"DELETE", "/teams/{id:[0-9]+}/repos/{owner:\\w+}/{repo:\\w+}"},
	{"GET", "/user/teams"},

	// Pull Requests
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/pulls"},
	//{"PATCH", "/repos/{owner}/{repo}/pulls/{number}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/commits"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/files"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/merge"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/merge"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/comments"},
	//{"GET", "/repos/{owner}/{repo}/pulls/comments"},
	//{"GET", "/repos/{owner}/{repo}/pulls/comments/{number}"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/pulls/{number:\\w+}/comments"},
	//{"PATCH", "/repos/{owner}/{repo}/pulls/comments/{number}"},
	//{"DELETE", "/repos/{owner}/{repo}/pulls/comments/{number}"},

	// Repositories
	{"GET", "/user/repos"},
	{"GET", "/users/{user:\\w+}/repos"},
	{"GET", "/orgs/{org:\\w+}/repos"},
	{"GET", "/repositories"},
	{"POST", "/user/repos"},
	{"POST", "/orgs/{org:\\w+}/repos"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}"},
	//{"PATCH", "/repos/{owner}/{repo}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/contributors"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/languages"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/teams"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/tags"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/branches"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/branches/{branch:\\w+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/collaborators"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/collaborators/{user:\\w+}"},
	{"PUT", "/repos/{owner:\\w+}/{repo:\\w+}/collaborators/{user:\\w+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/collaborators/{user:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/comments"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/commits/{sha:\\w+}/comments"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/commits/{sha:\\w+}/comments"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/comments/{id:[0-9]+}"},
	//{"PATCH", "/repos/{owner}/{repo}/comments/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/comments/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/commits"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/commits/{sha:\\w+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/readme"},
	//{"GET", "/repos/{owner}/{repo}/contents/*path"},
	//{"PUT", "/repos/{owner}/{repo}/contents/*path"},
	//{"DELETE", "/repos/{owner}/{repo}/contents/*path"},
	//{"GET", "/repos/{owner}/{repo}/:archive_format/:ref"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/keys"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/keys/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/keys"},
	//{"PATCH", "/repos/{owner}/{repo}/keys/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/keys/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/downloads"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/downloads/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/downloads/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/forks"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/forks"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/hooks"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/hooks"},
	//{"PATCH", "/repos/{owner}/{repo}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/hooks/{id:[0-9]+}/tests"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/hooks/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/merges"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/releases"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/releases/{id:[0-9]+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/releases"},
	//{"PATCH", "/repos/{owner:\\w+}/{repo}/releases/{id:[0-9]+}"},
	{"DELETE", "/repos/{owner:\\w+}/{repo:\\w+}/releases/{id:[0-9]+}"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/releases/{id:[0-9]+}/assets"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stats/contributors"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stats/commit_activity"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stats/code_frequency"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stats/participation"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/stats/punch_card"},
	{"GET", "/repos/{owner:\\w+}/{repo:\\w+}/statuses/{ref:\\w+}"},
	{"POST", "/repos/{owner:\\w+}/{repo:\\w+}/statuses/{ref:\\w+}"},

	// Search
	{"GET", "/search/repositories"},
	{"GET", "/search/code"},
	{"GET", "/search/issues"},
	{"GET", "/search/users"},
	{"GET", "/legacy/issues/search/{owner:\\w+}/{repo:\\w+}sitory/{state:\\w+}/{keyword:\\w+}"},
	{"GET", "/legacy/repos/search/{keyword:\\w+}"},
	{"GET", "/legacy/user/search/{keyword:\\w+}"},
	{"GET", "/legacy/user/email/{email:\\w+}"},

	// Users
	{"GET", "/users/{user:\\w+}"},
	{"GET", "/user"},
	//{"PATCH", "/user"},
	{"GET", "/users"},
	{"GET", "/user/emails"},
	{"POST", "/user/emails"},
	{"DELETE", "/user/emails"},
	{"GET", "/users/{user:\\w+}/followers"},
	{"GET", "/user/followers"},
	{"GET", "/users/{user:\\w+}/following"},
	{"GET", "/user/following"},
	{"GET", "/user/following/{user:\\w+}"},
	{"GET", "/users/{user:\\w+}/following/{target_user:\\w+}"},
	{"PUT", "/user/following/{user:\\w+}"},
	{"DELETE", "/user/following/{user:\\w+}"},
	{"GET", "/users/{user:\\w+}/keys"},
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
	beegoMuxRouter http.Handler
	boneRouter     http.Handler
	chiRouter      http.Handler
	httpRouter     http.Handler
	goRouter1      http.Handler
	goRouter2      http.Handler
	muxRouter      http.Handler
	trieMuxRouter  http.Handler
)

func init() {
	println("GithubAPI Routes:", len(githubAPI))
	println("GithubAPI2 Routes:", len(githubAPI2))

	calcMem("BeegoMuxRouter", func() {
		router := beegomux.New()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.method, route.path, handler)
		}
		beegoMuxRouter = router
	})

	calcMem("BoneRouter", func() {
		router := bone.New()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			if route.method == http.MethodGet {
				router.Get(route.path, http.HandlerFunc(handler))
			}
			if route.method == http.MethodPost {
				router.Post(route.path, http.HandlerFunc(handler))
			}
			if route.method == http.MethodPut {
				router.Put(route.path, http.HandlerFunc(handler))
			}
			if route.method == http.MethodDelete {
				router.Delete(route.path, http.HandlerFunc(handler))
			}
		}
		boneRouter = router
	})

	calcMem("ChiRouter", func() {
		router := chi.NewRouter()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.MethodFunc(route.method, route.path, handler)
		}
		chiRouter = router
	})

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

	calcMem("trie-mux", func() {
		router := triemux.New()
		handler := func(w http.ResponseWriter, _ *http.Request, _ triemux.Params) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.method, route.path, handler)
		}
		trieMuxRouter = router
	})

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

	calcMem("GoRouter1", func() {
		router := gorouter.New()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI {
			router.Handle(route.method, route.path, handler)
		}
		goRouter1 = router
	})

	calcMem("GoRouter2", func() {
		router := gorouter.New()
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(204)
		}
		for _, route := range githubAPI2 {
			router.Handle(route.method, route.path, handler)
		}
		goRouter2 = router
	})
}

// reference：https://github.com/julienschmidt/go-http-routing-benchmark/blob/2b136956a56bc65dddfa4bdaf7d1728ae2c90d50/bench_test.go#L76
func benchRoutes(b *testing.B, router http.Handler, routes []route) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, route := range routes {
			r.Method = route.method
			r.RequestURI = route.path
			u.Path = route.path
			u.RawQuery = rq
			router.ServeHTTP(w, r)
		}
	}
}

// With GithubAPI (goRouter vs beegoMuxRouter vs BoneRouter vs httpRouter vs trieMuxRouter)

func BenchmarkBeegoMuxRouterWithGithubAPI(b *testing.B) {
	benchRoutes(b, beegoMuxRouter, githubAPI)
}

func BenchmarkBoneRouterWithGithubAPI(b *testing.B) {
	benchRoutes(b, boneRouter, githubAPI)
}

func BenchmarkTrieMuxRouterWithGithubAPI(b *testing.B) {
	benchRoutes(b, trieMuxRouter, githubAPI)
}

func BenchmarkHttpRouterWithGithubAPI(b *testing.B) {
	benchRoutes(b, httpRouter, githubAPI)
}

func BenchmarkGoRouter1WithGithubAPI(b *testing.B) {
	benchRoutes(b, goRouter1, githubAPI)
}

// With GithubAPI2 (goRouter vs muxRouter vs chiRouter)

func BenchmarkGoRouter2WithGithubAPI2(b *testing.B) {
	benchRoutes(b, goRouter2, githubAPI2)
}

func BenchmarkChiRouterWithGithubAPI2(b *testing.B) {
	benchRoutes(b, chiRouter, githubAPI2)
}

func BenchmarkMuxRouterWithGithubAPI2(b *testing.B) {
	benchRoutes(b, muxRouter, githubAPI2)
}
