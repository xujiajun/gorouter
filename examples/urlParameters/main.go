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
		//get one URL parameter
		id := gorouter.GetParam(r, "id")
		//get all URL parameters
		//id := gorouter.GetAllParams(r)
		//fmt.Println(id)
		w.Write([]byte("match user/:id ! get id:" + id))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
