package gorouter_test

import (
	"fmt"
	"github.com/xujiajun/gorouter"
	//"log"
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

func TestRouter_CustomHandleNotFound(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/xxx", nil)

	if err != nil {
		t.Fatal(err)
	}

	customNotFoundStr := "404 page !"
	router.NotFoundFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, customNotFoundStr)
	})

	router.GET("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != customNotFoundStr {
		t.Errorf(errorFormat,
			rr.Body.String(), customNotFoundStr)
	}
}

func TestRouter_HandleNotFound(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/aaa", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/aa", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String()[:3] != "404" {
		t.Errorf(errorFormat,
			rr.Body.String(), "404 page not found\n")
	}
}

func TestGetParam(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	param := "1"
	req, err := http.NewRequest(http.MethodGet, "/test/"+param, nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/test/:id", func(w http.ResponseWriter, r *http.Request) {
		id := gorouter.GetParam(r, "id")
		if id != param {
			t.Fatal("TestGetParam test fail")
		}
	})
	router.ServeHTTP(rr, req)
}

func TestGetAllParams(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	param1 := "1"
	param2 := "2"
	req, err := http.NewRequest(http.MethodGet, "/param1/"+param1+"/param2/"+param2, nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param1/:param1/param2/:param2", func(w http.ResponseWriter, r *http.Request) {
		params := gorouter.GetAllParams(r)

		if params["param1"] != param1 {
			t.Fatal("TestGetAllParams test fail")
		}

		if params["param2"] != param2 {
			t.Fatal("TestGetAllParams test fail")
		}
	})
	router.ServeHTTP(rr, req)
}

func TestGetAllParamsMiss(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/param1", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param1", func(w http.ResponseWriter, r *http.Request) {
		params := gorouter.GetAllParams(r)

		if params != nil {
			t.Fatal("TestGetAllParams test fail")
		}

	})
	router.ServeHTTP(rr, req)
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Logged connection from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}

func TestRouter_Use(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.Use(withLogging)
	router.GET("/hi", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_UseForRoot(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.Use(withLogging)
	expected := "hi index"
	router.GET("/", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_Regex(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/param/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	router.GET("/param/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
		id := gorouter.GetParam(r, "id")
		if id != "1" {
			t.Fatal("TestGetAllParams test fail")
		}
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_HandleRoot(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	expected := "hi,root"

	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_HandlePanic(t *testing.T) {
	router := gorouter.New()

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered in f", r)
		}
	}()

	router.Handle("", "/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "invalid method")
	})

	router.ServeHTTP(rr, req)

	if rr.Body.String() != expected {
		t.Errorf(errorFormat,
			rr.Body.String(), expected)
	}
}

func TestRouter_Match(t *testing.T) {
	router := gorouter.New()
	requestUrl := "/xxx/1/yyy/2"

	ok := router.Match(requestUrl, "/xxx/:param1/yyy/:param2")

	if !ok {
		t.Fatal("TestRouter_Match test fail")
	}

	errorRequestUrl := "#xxx#1#yyy#2"
	ok = router.Match(errorRequestUrl, "/xxx/:param1/yyy/:param2")

	if ok {
		t.Fatal("TestRouter_Match test fail")
	}

	errorPath := "#xxx#1#yyy#2"
	ok = router.Match(requestUrl, errorPath)

	if ok {
		t.Fatal("TestRouter_Match test fail")
	}

	missRequestUrl := "/xxx/1/yyy/###"
	ok = router.Match(missRequestUrl, "/xxx/:param1/yyy/:param2")

	if ok {
		t.Fatal("TestRouter_Match test fail")
	}
}
