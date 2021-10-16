package myrouter

import (
	"net/http"
	"strings"
)

// tree is a trie tree.
type tree struct {
	node *node
}

type node struct {
	label    string
	actions  map[string]*action // key is method
	children map[string]*node   // key is a label o f next nodes
}

// action is an action.
type action struct {
	handler http.Handler
}

// result is a search result.
type result struct {
	actions *action
}

const (
	pathRoot      string = "/"
	pathDelimiter string = "/"
)

// newResult creates a new result.
func newResult() *result {
	return &result{}
}

// NewTree creates a new trie tree.
func NewTree() *tree {
	return &tree{
		node: &node{
			label:    pathRoot,
			actions:  make(map[string]*action),
			children: make(map[string]*node),
		},
	}
}

// Insert inserts a route definition to tree.
func (t *tree) Insert(methods []string, path string, handler http.Handler) error {
	curNode := t.node
	if path == pathRoot {
		curNode.label = path
		for _, method := range methods {
			curNode.actions[method] = &action{
				handler: handler,
			}
		}
		return nil
	}
	ep := explodePath(path)
	for i, p := range ep {
		nextNode, ok := curNode.children[p]
		if ok {
			curNode = nextNode
		}
		// Create a new node.
		if !ok {
			curNode.children[p] = &node{
				label:    p,
				actions:  make(map[string]*action),
				children: make(map[string]*node),
			}
			curNode = curNode.children[p]
		}
		// last loop.
		// If there is already registered data, overwrite it.
		if i == len(ep)-1 {
			curNode.label = p
			for _, method := range methods {
				curNode.actions[method] = &action{
					handler: handler,
				}
			}
			break
		}
	}

	return nil
}

func (t *tree) Search(method string, path string) (*result, error) {
	result := newResult()
	curNode := t.node
	if path != pathRoot {
		for _, p := range explodePath(path) {
			nextNode, ok := curNode.children[p]
			if !ok {
				if p == curNode.label {
					break
				} else {
					return nil, ErrNotFound
				}
			}
			curNode = nextNode
			continue
		}
	}
	result.actions = curNode.actions[method]
	if result.actions == nil {
		// no matching handler was found.
		return nil, ErrMethodNotAllowed
	}
	return result, nil
}

// explodePath removes an empty value in slice.
func explodePath(path string) []string {
	s := strings.Split(path, pathDelimiter)
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
