package main

import (
	"fmt"
	"net/http"

	myroute "github.com/bmf-san/introduction-to-golang-http-router-made-with-net-http"
)

func main() {
	r := myroute.NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodGet).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodGet).Handler(`/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/bar")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))

	http.ListenAndServe(":8080", r)
}
