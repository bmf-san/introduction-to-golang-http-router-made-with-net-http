package main

import (
	"fmt"
	"net/http"

	myroute "github.com/bmf-san/introduction-to-golang-http-router-made-with-net-http"
)

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /")
	})
}

func fooHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "GET /foo")
		case http.MethodPost:
			fmt.Fprintf(w, "POST /foo")
		default:
			fmt.Fprintf(w, "Not Found")
		}
	})
}

func fooBarHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /foo/bar")
	})
}

func fooBarBazHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /foo/bar/baz")
	})
}

func barHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /bar")
	})
}

func bazHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /baz")
	})
}

func main() {
	r := myroute.NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, indexHandler())
	r.Methods(http.MethodGet, http.MethodPost).Handler(`/foo`, fooHandler())
	r.Methods(http.MethodGet).Handler(`/foo/bar`, fooBarHandler())
	r.Methods(http.MethodGet).Handler(`/foo/bar/baz`, fooBarBazHandler())
	r.Methods(http.MethodGet).Handler(`/bar`, barHandler())
	r.Methods(http.MethodGet).Handler(`/baz`, bazHandler())

	http.ListenAndServe(":8080", r)
}
