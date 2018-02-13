package main

import (
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

func main() {
	mux := gorouter.New()
	//url parameters match
	mux.GET("/user/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("match user/:id !"))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
