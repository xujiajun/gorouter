package gorouter_test

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	errorFormat, expected string
)

func init() {
	expected = "hi,gorouter"
	errorFormat = "handler returned unexpected body: got %v want %v"
}

func TestRouter_GET(t *testing.T) {

	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_POST(t *testing.T) {

	router := gorouter.New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.POST("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_DELETE(t *testing.T) {

	router := gorouter.New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodDelete, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.DELETE("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_PATCH(t *testing.T) {

	router := gorouter.New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPatch, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.PATCH("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_PUT(t *testing.T) {

	router := gorouter.New()
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPut, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.PUT("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_Group(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	prefix := "/api"

	req, err := http.NewRequest(http.MethodGet, prefix+"/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.Group(prefix).GET("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}
