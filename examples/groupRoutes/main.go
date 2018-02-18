package main

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "/api/users")
}

func main() {
	mux := gorouter.New()
	mux.Group("/api").GET("/users", usersHandler)

	log.Fatal(http.ListenAndServe(":8181", mux))
}
