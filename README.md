# gorouter  [![GoDoc](https://godoc.org/github.com/xujiajun/gorouter?status.svg)](https://godoc.org/github.com/xujiajun/gorouter) <a href="https://travis-ci.org/xujiajun/gorouter"><img src="https://travis-ci.org/xujiajun/gorouter.svg?branch=master" alt="Build Status"></a> [![Go Report Card](https://goreportcard.com/badge/github.com/xujiajun/gorouter)](https://goreportcard.com/report/github.com/xujiajun/gorouter)
A simple and fast HTTP router for Go.

## Motivation

I wanted a simple, fast router that has no unnecessary overhead using the standard library only, following good practices and well tested code.

## Features

* Fast
* URL parameters
* Regex parameters
* Routes groups
* Custom NotFoundHandler
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

### Regex Parameters

### Group routes

### Custom NotFoundHandler

### Middlewares

## Contributing

If you'd like to help out with the project. You can put up a Pull Request.

## License

The gorouter is open-sourced software licensed under the [MIT Licensed](http://www.opensource.org/licenses/MIT)
