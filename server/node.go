package server

import "net/http"

type matchFunc func(path string, c *Context) bool

type Node struct {
	children []*Node
	pattern  string
	handler  HandlerFunc
	isMatch  matchFunc
}

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

/** create a root tree **/
func CreateRootNode(method string) *Node {
	return &Node{
		children: make([]*Node, 0, 2),
		pattern:  method,
		isMatch: func(path string, c *Context) bool {
			panic("cannot call this at root")
		},
	}
}

/** create a method tree **/
func CreateMethodTree() map[string]*Node {
	forest := make(map[string]*Node, len(supportMethods))
	for _, m := range supportMethods {
		forest[m] = CreateRootNode(m)
	}
	return forest
}

/** create a node **/
func CreateNode(path string) *Node {
	return &Node{
		children: make([]*Node, 0, 2),
		pattern:  path,
		isMatch: func(p string, c *Context) bool {
			return p == path && p != "*"
		},
	}
}

/** create a sub node **/
func CreateSubNode(node *Node, paths []string, handlerFunc HandlerFunc) {
	cur := node
	for _, path := range paths {
		newNode := CreateNode(path)
		cur.children = append(cur.children, newNode)
		cur = newNode
	}
	cur.handler = handlerFunc
}
