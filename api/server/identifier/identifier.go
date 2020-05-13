package identifier

import (
	"bennu.cl/identifier-producer/api/server"
	"bennu.cl/identifier-producer/pkg/identifier"
)

type identifierRouter struct {
	api    identifier.Service
	routes []server.Route
}

func NewRouter(ids identifier.Service) server.Router {
	r := &identifierRouter{
		api: ids,
	}
	r.initRoutes()
	return r
}

func (idr *identifierRouter) Routes() []server.Route {
	return idr.routes
}

func (idr *identifierRouter) initRoutes() {
	idr.routes = []server.Route{
		server.NewPostRoute("/identifier", idr.postIdentifier),
	}
}
