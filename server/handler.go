package server

import (
	"net/http"
)

type Handler interface {
	ServeHTTP(c *Context)
}

type HandlerFunc func(c *Context)

type WebHandler struct {
	router Router
}

/** default handler **/
func (h *WebHandler) ServeHTTP(c *Context) {
	/** find route **/
	handler, found := h.router.Find(c.R.Method, c.R.URL.Path, c)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)
}

/** create a handler instance **/	
func NewHandler() *WebHandler {
	return &WebHandler{
		router: NewRouter(),
	}
}
