package server

import (
	"net/http"
)

type Server interface {
	Start(address string) error
	Route(method string, pattern string, handlerFn HandlerFunc) error
}

type WebServer struct {
	handler WebHandler
	fnc     HandlerFunc
}

/** start server **/
func (s *WebServer) Start(address string) error {
	return http.ListenAndServe(address, s)
}

/** register route **/
func (s *WebServer) Route(method string, pattern string, handlerFn HandlerFunc) error {
	return s.handler.router.Create(method, pattern, handlerFn)
}

/** default handler **/
func (s *WebServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := NewContext(writer, request)
	s.fnc(c)
}

/** create a server instance **/
func NewServer(middleWareList ...MiddleWare) Server {
	/** create a handler instance **/	
	newHandler := NewHandler()
	fnc := newHandler.ServeHTTP

	/** handle func **/	
	for i := len(middleWareList) - 1; i >= 0; i-- {
		m := middleWareList[i]
		fnc = m(fnc)
	}

	return &WebServer{
		handler: *newHandler,
		fnc:     fnc,
	}
}
