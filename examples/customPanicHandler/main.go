package main

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

func main() {
	mux := gorouter.New()
	mux.PanicHandler = func(w http.ResponseWriter, req *http.Request, err interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("err from recover is :", err)
		fmt.Fprint(w, "received a panic")
	}
	mux.GET("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
