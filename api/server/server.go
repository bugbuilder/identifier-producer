package server

import (
	"context"
	"github.com/gorilla/mux"
	"k8s.io/klog"
	"net/http"
	"time"
)

type Server struct {
	routers []Router
	server  *http.Server
}

func NewServer() *Server {
	return &Server{
		server: &http.Server{
			Addr:              ":8080",
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      300 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

func (s *Server) InitRouter(routers ...Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) makeHTTPHandler(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc := handler

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(w, r, vars); err != nil {
			statusCode := getHTTPErrorStatusCode(err)
			if statusCode >= 500 {
				klog.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			}
			MakeErrorHandler(statusCode, err)(w, r)
		}
	}
}

func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()

	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	return m
}

func (s *Server) run() error {
	errRun := make(chan error)
	s.server.Handler = s.createMux()
	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errRun <- err
		}
	}(s.server)

	err := <-errRun
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run(errRun chan error) {
	if err := s.run(); err != nil {
		errRun <- err
		return
	}
	errRun <- nil
}

func (s *Server) Shutdown(errRun chan error) {
	klog.Infof("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		klog.Fatal("server forced to shutdown: %s", err)
		errRun <- err
	}

	errRun <- nil
}
