package myrouter

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewResult(t *testing.T) {
	actual := newResult()
	expected := &result{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

func TestNewTree(t *testing.T) {
	actual := NewTree()
	expected := &tree{
		node: &node{
			label:    pathRoot,
			actions:  make(map[string]*action),
			children: make(map[string]*node),
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

// item is a set of routing definition.
type item struct {
	method string
	path   string
}

// caseWithFailure is a struct for testWithFailure.
type caseWithFailure struct {
	hasError bool
	item     *item
	expected *result
}

func TestInsert(t *testing.T) {
	tree := NewTree()
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	cases := []struct {
		method  string
		path    string
		handler http.Handler
	}{
		{
			method:  http.MethodGet,
			path:    "/",
			handler: fooHandler,
		},
		{
			method:  http.MethodPost,
			path:    "/",
			handler: fooHandler,
		},
		{
			method:  http.MethodGet,
			path:    "/foo",
			handler: fooHandler,
		},
		{
			method:  http.MethodPost,
			path:    "/foo",
			handler: fooHandler,
		},
		{
			method:  http.MethodGet,
			path:    "/foo/bar",
			handler: fooHandler,
		},
		{
			method:  http.MethodPost,
			path:    "/foo/bar",
			handler: fooHandler,
		},
		{
			method:  http.MethodGet,
			path:    "/foo/bar/baz",
			handler: fooHandler,
		},
		{
			method:  http.MethodPost,
			path:    "/foo/bar/baz",
			handler: fooHandler,
		},
	}

	for _, c := range cases {
		if err := tree.Insert([]string{c.method}, c.path, c.handler); err != nil {
			t.Errorf("err: %v\n", err)
		}
	}
}

func TestSearch(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler)
	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler)
	tree.Insert([]string{http.MethodGet}, `/bar`, barHandler)
	tree.Insert([]string{http.MethodGet}, `/foo/bar`, fooBarHandler)

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler: rootHandler,
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler: fooHandler,
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &result{
				actions: &action{
					handler: barHandler,
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/baz",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler: fooBarHandler,
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/baz",
			},
			expected: nil,
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/baz",
			},
			expected: nil,
		},
	}

	testWithFailure(t, tree, cases)
}

func testWithFailure(t *testing.T, tree *tree, cases []caseWithFailure) {
	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)

		if c.hasError {
			if err == nil {
				t.Fatalf("actual: %v expected err: %v", actual, err)
			}

			if actual != c.expected {
				t.Errorf("actual:%v expected:%v", actual, c.expected)
			}

			continue
		}

		if err != nil {
			t.Fatalf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.actions.handler) != reflect.ValueOf(c.expected.actions.handler) {
			t.Errorf("actual:%v expected:%v", actual.actions.handler, c.expected.actions.handler)
		}
	}
}

func TestExplodePath(t *testing.T) {
	cases := []struct {
		actual   []string
		expected []string
	}{
		{
			actual:   explodePath(""),
			expected: nil,
		},
		{
			actual:   explodePath("/"),
			expected: nil,
		},
		{
			actual:   explodePath("//"),
			expected: nil,
		},
		{
			actual:   explodePath("///"),
			expected: nil,
		},
		{
			actual:   explodePath("/foo"),
			expected: []string{"foo"},
		},
		{
			actual:   explodePath("/foo/bar"),
			expected: []string{"foo", "bar"},
		},
		{
			actual:   explodePath("/foo/bar/baz"),
			expected: []string{"foo", "bar", "baz"},
		},
		{
			actual:   explodePath("/foo/bar/baz/"),
			expected: []string{"foo", "bar", "baz"},
		},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
