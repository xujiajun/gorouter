# gorouter  [![GoDoc](https://godoc.org/github.com/xujiajun/gorouter?status.svg)](https://godoc.org/github.com/xujiajun/gorouter) <a href="https://travis-ci.org/xujiajun/gorouter"><img src="https://travis-ci.org/xujiajun/gorouter.svg?branch=master" alt="Build Status"></a> [![Go Report Card](https://goreportcard.com/badge/github.com/xujiajun/gorouter)](https://goreportcard.com/report/github.com/xujiajun/gorouter) [![Coverage Status](https://s3.amazonaws.com/assets.coveralls.io/badges/coveralls_100.svg)](https://coveralls.io/github/xujiajun/gorouter?branch=master) [![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xujiajun/gorouter/master/LICENSE)  [![Release](https://img.shields.io/badge/release-v1.1.0-blue.svg?style=flat-square)](https://github.com/xujiajun/gorouter/releases/tag/v1.0.1) [![Awesome](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#routers) 
`xujiajun/gorouter` is a simple and fast HTTP router for Go. It is easy to build RESTful APIs and your web framework.

## Motivation

I wanted a simple and fast HTTP GO router, which supports regexp. I prefer to support regexp is because otherwise it will need the logic to check the URL parameter type, thus increasing the program complexity. So I did some searching on Github and found the wonderful `julienschmidt/httprouter`: it is very fast，unfortunately it does not support regexp. Later I found out about `gorilla/mux`: it is powerful as well，but a written benchmark shows me that it is somewhat slow. So I tried to develop a new router which both supports regexp and should be fast. Finally I did it and named `xujiajun/gorouter`. By the way, this is my first GO open source project. It may be the fastest GO HTTP router which supports regexp, and regarding its performance please refer to my latest [Benchmarks](#benchmarks).

## Features

* Fast - see [Benchmarks](#benchmarks)
* [URL parameters](#url-parameters)
* [Regex parameters](#regex-parameters)
* [Routes groups](#routes-groups)
* [Reverse Routing](#reverse-routing)
* [Custom NotFoundHandler](#custom-notfoundhandler)
* [Custom PanicHandler](#custom-panichandler)
* [Middleware Chain Support](#middlewares-chain)
* [Serve Static Files](#serve-static-files)
* [Pattern Rule Familiar](#pattern-rule)
* HTTP Method Get、Post、Delete、Put、Patch Support
* No external dependencies (just Go stdlib)


## Requirements

* golang 1.8+

## Installation

```
go get -u github.com/xujiajun/gorouter
```

## Usage

### Static routes

```golang
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

```golang
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
```

### Regex Parameters

```golang
package main

import (
	"log"
	"net/http"

	"github.com/xujiajun/gorouter"
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


### Routes Groups

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xujiajun/gorouter"
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

### Reverse Routing

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/xujiajun/gorouter"
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

```

### Custom NotFoundHandler

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xujiajun/gorouter"
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

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xujiajun/gorouter"
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

### Middlewares Chain

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xujiajun/gorouter"
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

```golang
package main

import (
	"log"
	"net/http"
	"os"
	
	"github.com/xujiajun/gorouter"
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

> if use default regex checks unless you know what you're doing

 
## Benchmarks

> go test -bench=.

### Benchmark System:

* Go Version : go1.11.2 darwin/amd64
* OS:        Mac OS X 10.13.6
* Architecture:   x86_64
* 16 GB 2133 MHz LPDDR3
* CPU: 3.1 GHz Intel Core i7

### Tested routers:

* [beego/mux](https://github.com/beego/mux)
* [go-zoo/bone](https://github.com/go-zoo/bone)
* [go-chi/chi](https://github.com/go-chi/chi)
* [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
* [gorilla/mux](https://github.com/gorilla/mux)
* [trie-mux/mux](https://github.com/teambition/trie-mux)
* [xujiajun/gorouter](https://github.com/xujiajun/gorouter)


Thanks the author of httprouter: [@julienschmidt](https://github.com/julienschmidt) give me advise about benchmark [issues/24](https://github.com/xujiajun/gorouter/issues/24)

## Result:

Given some routing matching syntax differences, divide GithubAPI into two groups：

### Using GithubAPI Result：

```
BenchmarkBeegoMuxRouterWithGithubAPI-8   	   10000	    142398 ns/op	  134752 B/op	    1038 allocs/op
BenchmarkBoneRouterWithGithubAPI-8       	    1000	   2104486 ns/op	  720160 B/op	    8620 allocs/op
BenchmarkTrieMuxRouterWithGithubAPI-8    	   20000	     80845 ns/op	   65856 B/op	     537 allocs/op
BenchmarkHttpRouterWithGithubAPI-8       	   50000	     30169 ns/op	   13792 B/op	     167 allocs/op
BenchmarkGoRouter1WithGithubAPI-8        	   30000	     57793 ns/op	   13832 B/op	     406 allocs/op

```
### Using GithubAPI2 Result：

```
BenchmarkGoRouter2WithGithubAPI2-8       	   30000	     57613 ns/op	   13832 B/op	     406 allocs/op
BenchmarkChiRouterWithGithubAPI2-8       	   10000	    143224 ns/op	  104436 B/op	    1110 allocs/op
BenchmarkMuxRouterWithGithubAPI2-8       	     300	   4450731 ns/op	   61463 B/op	     995 allocs/op
```

### All togther Result：

```
➜  gorouter git:(master) go test -bench=.
GithubAPI Routes: 203
GithubAPI2 Routes: 203
   BeegoMuxRouter: 111072 Bytes
   BoneRouter: 100992 Bytes
   ChiRouter: 71512 Bytes
   HttpRouter: 37016 Bytes
   trie-mux: 131128 Bytes
   MuxRouter: 1378496 Bytes
   GoRouter1: 83824 Bytes
   GoRouter2: 85584 Bytes
goos: darwin
goarch: amd64
pkg: github.com/xujiajun/gorouter
BenchmarkBeegoMuxRouterWithGithubAPI-8   	   10000	    142398 ns/op	  134752 B/op	    1038 allocs/op
BenchmarkBoneRouterWithGithubAPI-8       	    1000	   2104486 ns/op	  720160 B/op	    8620 allocs/op
BenchmarkTrieMuxRouterWithGithubAPI-8    	   20000	     80845 ns/op	   65856 B/op	     537 allocs/op
BenchmarkHttpRouterWithGithubAPI-8       	   50000	     30169 ns/op	   13792 B/op	     167 allocs/op
BenchmarkGoRouter1WithGithubAPI-8        	   30000	     57793 ns/op	   13832 B/op	     406 allocs/op
BenchmarkGoRouter2WithGithubAPI2-8       	   30000	     57613 ns/op	   13832 B/op	     406 allocs/op
BenchmarkChiRouterWithGithubAPI2-8       	   10000	    143224 ns/op	  104436 B/op	    1110 allocs/op
BenchmarkMuxRouterWithGithubAPI2-8       	     300	   4450731 ns/op	   61463 B/op	     995 allocs/op
PASS
ok  	github.com/xujiajun/gorouter	15.918s

```

### Conclusions:

* Performance (xujiajun/gorouter,julienschmidt/httprouter and teambition/trie-mux is fast)

* Memory Consumption (xujiajun/gorouter and julienschmidt/httprouter is fewer) 

* Features (julienschmidt/httprouter not support regexp，but others support it)

> if you want a high performance router which support regexp, maybe [xujiajun/gorouter](https://github.com/xujiajun/gorouter) is good choice.

> if you want a high performance router which not support regexp, maybe [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) is good choice.

In the end, as julienschmidt said `performance can not be the (only) criterion for choosing a router. Play around a bit with some of the routers, and choose the one you like best.`

## Contributing

If you'd like to help out with the project. You can put up a Pull Request. Thanks to all [contributors](https://github.com/xujiajun/gorouter/graphs/contributors).

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
