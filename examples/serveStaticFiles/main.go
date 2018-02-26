package main

import (
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
	"os"
)

//ServeFiles serve static resources
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir := wd + "/examples/serveStaticFiles/files"
	http.StripPrefix("/files/", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
}

func main() {
	mux := gorouter.New()
	mux.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	//defined prefix
	mux2 := mux.Group("/files")
	//http://127.0.0.1:8181/files/demo.txt
	//will match
	mux2.GET("/{filename:[0-9a-zA-Z_.]+}", func(w http.ResponseWriter, r *http.Request) {
		ServeFiles(w, r)
	})

	//http://127.0.0.1:8181/files/a/demo2.txt
	//http://127.0.0.1:8181/files/a/demo.txt
	//will match
	mux2.GET("/{fileDir:[0-9a-zA-Z_.]+}/{filename:[0-9a-zA-Z_.]+}", func(w http.ResponseWriter, r *http.Request) {
		ServeFiles(w, r)
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
