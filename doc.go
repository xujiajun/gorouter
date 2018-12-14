// Copyright 2018 The xujiajun/gorouter Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package xujiajun/gorouter is a simple and fast HTTP router for Go. It is easy to build RESTful APIs and your web framework.
//
// Here is the example:
//
//  package main
//
//  import (
//	 "log"
//	 "net/http"
//
//	 "github.com/xujiajun/gorouter"
//  )
//
//  func main() {
//	 mux := gorouter.New()
//	 //url parameters match
//	 mux.GET("/user/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
//		//get one URL parameter
//		id := gorouter.GetParam(r, "id")
//		//get all URL parameters
//		//id := gorouter.GetAllParams(r)
//		//fmt.Println(id)
//		w.Write([]byte("match user/{id:[0-9]+} ! get id:" + id))
//	  })
//
//	  log.Fatal(http.ListenAndServe(":8181", mux))
//  }
//
// Here is the syntax:
//
//  Syntax	                   Description	                          Example
//  :name	            named parameter	                        /user/:name
//  {name:regexp}	    named with regexp parameter	            /user/{name:[0-9a-zA-Z]+}
//  :id	                named with regexp parameter	            /user/:id
//
//  And :id is short for {id:[0-9]+}, :name are short for {name:[0-9a-zA-Z_]+}.
package gorouter
