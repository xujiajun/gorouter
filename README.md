# gorouter  [![GoDoc](https://godoc.org/github.com/xujiajun/gorouter?status.svg)](https://godoc.org/github.com/xujiajun/gorouter) <a href="https://travis-ci.org/xujiajun/gorouter"><img src="https://travis-ci.org/xujiajun/gorouter.svg?branch=master" alt="Build Status"></a> [![Go Report Card](https://goreportcard.com/badge/github.com/xujiajun/gorouter)](https://goreportcard.com/report/github.com/xujiajun/gorouter) [![Coverage Status](https://s3.amazonaws.com/assets.coveralls.io/badges/coveralls_100.svg)](https://coveralls.io/github/xujiajun/gorouter?branch=master) [![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xujiajun/gorouter/master/LICENSE)
A simple and fast HTTP router for Go.

## Motivation

I wanted a simple, fast router that has no unnecessary overhead using the standard library only, following good practices and well tested code.

## Features

* Fast - see [benchmarks](#benchmarks)
* URL parameters
* Regex parameters
* Routes groups
* Custom NotFoundHandler
* Custom PanicHandler
* Middleware chain Support
* No external dependencies (just Go 1.7+ stdlib)


## Installation

```
go get github.com/xujiajun/gorouter
```

## Usage

### Static routes

```
package main

import (
	"log"
	"net/http"
	"github.com/xujiajun/gorouter"
)

func main() {
	mux := gorouter.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	log.Fatal(http.ListenAndServe(":8181", mux))
}

```

### URL Parameters

```
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
```

### Regex Parameters

```
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
```


### Group routes

```
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
```

### Custom NotFoundHandler

```
package main

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 page !!!")
}

func main() {
	mux := gorouter.New()
	mux.NotFoundFunc(notFoundFunc)
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
```

### Custom PanicHandler

```
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

```

### Middlewares

```
package main

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"log"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

//https://upgear.io/blog/golang-tip-wrapping-http-response-writer-for-middleware/
func withStatusRecord(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, http.StatusOK}
		next.ServeHTTP(&rec, r)
		log.Printf("response status: %v\n", rec.status)
	}
}

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not found page !")
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func withTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func main() {
	mux := gorouter.New()
	mux.NotFoundFunc(notFoundFunc)
	mux.Use(withLogging, withTracing, withStatusRecord)
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	log.Fatal(http.ListenAndServe(":8181", mux))
}
```

## Serve static files

```
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
```
Detail see [serveStaticFiles example](https://github.com/xujiajun/gorouter/blob/master/examples/serveStaticFiles/main.go)

## Pattern Rule

The syntax here is modeled after [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) and [gorilla/mux](https://github.com/gorilla/mux)

| Syntax | Description | Example |
|--------|------|-------|
| `:name` | named parameter | /user/:name |
| `{name:regexp}` | named with regexp parameter |  /user/{name:[0-9a-zA-Z]+} |
| `:id` | named with regexp parameter |  /user/:id |

And `:id` is short for `{id:[0-9]+}`, `:name` are short for `{name:[0-9a-zA-Z_]+}`

 
## Benchmarks

> go test -bench=.

Benchmark System:

* Go version 1.9.2
* OS:        Mac OS X 10.13.3 
* Architecture:   x86_64
* 16 GB 2133 MHz LPDDR3

Tested routers:

* [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
* [xujiajun/GoRouter](https://github.com/xujiajun/gorouter)
* [gorilla/mux](https://github.com/gorilla/mux)


Result:

```
GithubAPI Routes: 203
GithubAPI2 Routes: 203
   HttpRouter: 37464 Bytes
   GoRouter: 83616 Bytes
   MuxRouter: 1324192 Bytes
goos: darwin
goarch: amd64
pkg: github.com/xujiajun/gorouter
BenchmarkHttpRouter-8   	   10000	    517116 ns/op	 1034339 B/op	    2604 allocs/op
BenchmarkGoRouter-8     	   10000	    551223 ns/op	 1034385 B/op	    2843 allocs/op
BenchmarkMuxRouter-8    	   10000	   5825422 ns/op	 1272958 B/op	    4691 allocs/op
PASS
ok  	github.com/xujiajun/gorouter	68.965s
```

Conclusions:

* Memory Consumption (HttpRouter > gorouter > MuxRouter) 

* Performance (HttpRouter > gorouter > MuxRouter)

* Features (HttpRouter not support regexp, But GoRouter and MuxRouter support)

As author of [HttpRouter](https://github.com/julienschmidt/httprouter) said `performance can not be the (only) criterion for choosing a router. Play around a bit with some of the routers, and choose the one you like best. Moreover main memory is cheap and usually not a scarce resource.`


## Contributing

If you'd like to help out with the project. You can put up a Pull Request.

## Author

* [xujiajun](https://github.com/xujiajun)

## License

The gorouter is open-sourced software licensed under the [MIT Licensed](http://www.opensource.org/licenses/MIT)

## Acknowledgements

This package is inspired by the following:

* [httprouter](https://github.com/julienschmidt/httprouter)
* [bone](https://github.com/go-zoo/bone)
* [trie-mux](https://github.com/teambition/trie-mux)
* [alien](https://github.com/gernest/alien)
