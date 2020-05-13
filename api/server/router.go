package server

import "net/http"

type Router interface {
	Routes() []Route
}

type Route interface {
	Handler() APIFunc
	Method() string
	Path() string
}

type APIFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) error

type route struct {
	method  string
	path    string
	handler APIFunc
}

func (idr route) Handler() APIFunc {
	return idr.handler
}

func (idr route) Method() string {
	return idr.method
}

func (idr route) Path() string {
	return idr.path
}

func NewRoute(method, path string, handler APIFunc) Route {
	return route{method, path, handler}
}

func NewPostRoute(path string, handler APIFunc) Route {
	return NewRoute(http.MethodPost, path, handler)
}
