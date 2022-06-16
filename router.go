package myrouter

import (
	"errors"
	"net/http"
)

var (
	Http404Response = []byte("page not found")
	Http405Response = []byte("method not allowed")
)

func Set404(content string) {
	Http404Response = []byte(content)
}

func Set405(content string) {
	Http404Response = []byte(content)
}

// Router represents the router which handles routing.
type Router struct {
	tree *tree
}

// route represents the route which has data for a routing.
type route struct {
	methods []string
	path    string
	handler http.Handler
}

var (
	tmpRoute = &route{}
	// Error for not found.
	ErrNotFound = errors.New("no matching route was found")
	// Error for method not allowed.
	ErrMethodNotAllowed = errors.New("methods is not allowed")
)

func NewRouter() *Router {
	return &Router{
		tree: NewTree(),
	}
}

func (r *Router) Methods(methods ...string) *Router {
	tmpRoute.methods = append(tmpRoute.methods, methods...)
	return r
}

// Handler sets a handler.
func (r *Router) Handler(path string, handler http.Handler) {
	tmpRoute.handler = handler
	tmpRoute.path = path
	r.Handle()
}

// Handle handles a route.
func (r *Router) Handle() {
	r.tree.Insert(tmpRoute.methods, tmpRoute.path, tmpRoute.handler)
	tmpRoute = &route{}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	result, err := r.tree.Search(method, path)
	if err != nil {
		status, msg := handleErr(err)
		w.WriteHeader(status)
		w.Write(msg)
		return
	}
	h := result.actions.handler
	h.ServeHTTP(w, req)
}

func handleErr(err error) (int, []byte) {
	var status int
	var body []byte
	switch err {
	case ErrMethodNotAllowed:
		status = http.StatusMethodNotAllowed
		body = Http405Response
	case ErrNotFound:
		status = http.StatusNotFound
		body = Http404Response
	}
	return status, body
}
