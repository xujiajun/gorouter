package main

import (
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

func main() {
	mux := gorouter.New()
	//url regex match
	mux.GET("/user/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("match user/{id:[0-9]+} !"))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
