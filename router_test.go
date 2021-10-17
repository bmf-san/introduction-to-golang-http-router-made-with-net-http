package myrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	actual := NewRouter()
	expected := &Router{
		tree: NewTree(),
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v\n", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodGet).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.Methods(http.MethodPost).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodPost).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodPost).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.Methods(http.MethodOptions).Handler(`/options`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/options")
	}))

	cases := []struct {
		path   string
		method string
		code   int
		body   string
	}{
		{
			path:   "/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			path:   "/foo",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo",
		},
		{
			path:   "/foo/bar",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar",
		},
		{
			path:   "/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			path:   "/foo",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo",
		},
		{
			path:   "/foo/bar",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar",
		},
		{
			path:   "/options",
			method: http.MethodOptions,
			code:   http.StatusOK,
			body:   "/options",
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest(c.method, c.path, nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != c.code {
			t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
		}

		recBody, _ := ioutil.ReadAll(rec.Body)
		body := string(recBody)
		if body != c.body {
			t.Errorf("actual: %v expected: %v\n", body, c.body)
		}
	}
}

func TestHandleErr(t *testing.T) {
	cases := []struct {
		actual   int
		expected int
	}{
		{
			actual:   handleErr(ErrMethodNotAllowed),
			expected: http.StatusMethodNotAllowed,
		},
		{
			actual:   handleErr(ErrNotFound),
			expected: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual: %v expected: %v\n", c.actual, c.expected)
		}
	}

}
