package server

import (
	"fmt"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) Respond(content string) {
	fmt.Fprintf(c.W, "%s", content)
}
func NewContext(W http.ResponseWriter, R *http.Request) *Context {
	return &Context{
		W: W,
		R: R,
	}
}
