package server

import (
	"errors"
	"strings"
)

type Router interface {
	Create(method string, pattern string, handler HandlerFunc) error
	Find(method string, path string, c *Context) (HandlerFunc, bool)
}

type WebRouter struct {
	methodTree map[string]*Node
}

/** register a route **/
func (r *WebRouter) Create(method string, pattern string, handler HandlerFunc) error {
	err := r.validatePattern(pattern)
	if err != nil {
		return err
	}

	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur, ok := r.methodTree[method]
	if !ok {
		return errors.New("invalid method")
	}

	for index, path := range paths {
		matchChild, found := r.findMatchChild(cur, path, nil)
		if found {
			cur = matchChild
		} else {
			CreateSubNode(cur, paths[index:], handler)
			return nil
		}
	}
	cur.handler = handler
	return nil
}

/** find  route **/
func (r *WebRouter) Find(method string, path string, c *Context) (HandlerFunc, bool) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur, ok := r.methodTree[method]
	if !ok {
		return nil, false
	}
	for _, p := range paths {
		matchChild, found := r.findMatchChild(cur, p, c)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	if cur.handler == nil {
		// 到达这里是因为这种场景
		// 比如说你注册了 /user/profile
		// 然后你访问 /user
		return nil, false
	}
	return cur.handler, true
}

/** validate route pattern **/
func (r *WebRouter) validatePattern(pattern string) error {
	pos := strings.Index(pattern, "*")
	if pos > 0 {
		if pos != len(pattern)-1 {
			return errors.New("invalid router pattern")
		}
		if pattern[pos-1] != '/' {
			return errors.New("invalid router pattern")
		}
	}
	return nil
}

/** find sub route in tree **/
func (r *WebRouter) findMatchChild(node *Node, path string, c *Context) (*Node, bool) {
	candidates := make([]*Node, 0, 2)
	for _, child := range node.children {
		if child.isMatch(path, c) {
			candidates = append(candidates, child)
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	// sort.Slice(candidates, func(i, j int) bool {
	// 	return candidates[i].nodeType < candidates[j].nodeType
	// })
	return candidates[len(candidates)-1], true
}

/** create a router instance **/
func NewRouter() Router {
	return &WebRouter{
		methodTree: CreateMethodTree(),
	}
}
