package server

import "net/http"

type matchFunc func(path string, c *Context) bool

type Node struct {
	children []*Node
	pattern  string
	handler  HandlerFunc
	isMatch  matchFunc
	nodeType int
}

/** define the sort of nodeType **/
const (
	nodeTypeRoot = iota    
	nodeTypeAny
	nodeTypeNormal
)

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
		nodeType: nodeTypeRoot,
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
func createNode(path string) *Node {
	if path == "*"{
		return createAnyNode()
	}else{
		return createNormalNode(path)
	}
}

/** create a normal node **/
func createNormalNode(path string) *Node {
	return &Node{
		children: make([]*Node, 0, 2),
		pattern:  path,
		isMatch: func(p string, c *Context) bool {
			return p == path && p != "*"
		},
		nodeType:nodeTypeNormal,
	}
}

/** create an any node **/
func createAnyNode() *Node{
	/** any node does not have children **/
	return &Node{
		isMatch: func(p string, c *Context) bool {
			return true
		},
		pattern: "*",
		nodeType: nodeTypeAny,
	}
}

/** create a sub node **/
func CreateSubNode(node *Node, paths []string, handlerFunc HandlerFunc) {
	cur := node
	for _, path := range paths {
		newNode := createNode(path)
		cur.children = append(cur.children, newNode)
		cur = newNode
	}
	cur.handler = handlerFunc
}
