package main

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"net/http"
)

func main() {
	mux := gorouter.New()

	routeName1 := "user_event"
	mux.GETAndName("/users/:user/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/:user/events"))
	}, routeName1)

	routeName2 := "repos_owner"
	mux.GETAndName("/repos/{owner:\\w+}/{repo:\\w+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/repos/{owner:\\w+}/{repo:\\w+}"))
	}, routeName2)

	params := make(map[string]string)
	params["user"] = "xujiajun"
	fmt.Println(mux.Generate(http.MethodGet, routeName1, params)) // /users/xujiajun/events <nil>

	params = make(map[string]string)
	params["owner"] = "xujiajun"
	params["repo"] = "xujiajun_repo"
	fmt.Println(mux.Generate(http.MethodGet, routeName2, params)) // /repos/xujiajun/xujiajun_repo <nil>
}
