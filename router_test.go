package gorouter_test

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func hiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hi,gorouter")
}

func TestRouter_GET(t *testing.T) {

	req, err := http.NewRequest("GET", "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := gorouter.New()

	router.GET("/hi", hiHandler)
	router.ServeHTTP(rr, req)

	expected := "hi,gorouter"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
