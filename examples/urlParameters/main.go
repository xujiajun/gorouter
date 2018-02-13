package main

import (
	"log"
	"net/http"
	"github.com/xujiajun/gorouter"
)

func main() {
	mux := gorouter.New()
	//url parameters match
	mux.GET("/user/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("match user/:id !"))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
