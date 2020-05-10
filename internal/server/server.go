package server

import (
	"bennu.cl/identifier-producer/internal/handlers"
	"bennu.cl/identifier-producer/pkg/api"
	"bennu.cl/identifier-producer/pkg/kafka"
	"context"
	"k8s.io/klog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	identifier api.IdentifierService
	healthz    kafka.Healthz
}

func NewServer(ids api.IdentifierService, healthz kafka.Healthz) *Server {
	srv := &Server{}

	srv.identifier = ids
	srv.healthz = healthz

	return srv
}

func (srv *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/identifier", handlers.Producer(srv.identifier))
	mux.HandleFunc("/healthz", handlers.Healthz(srv.healthz))

	s := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	klog.Infof("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exitCode := 0
	if err := s.Shutdown(ctx); err != nil {
		klog.Fatal("Server forced to shutdown: %s", err)
		exitCode = 1
	}

	klog.Infof("awaiting kafka shutdown")
	time.Sleep(1 * time.Second)

	klog.Infof("Exiting with %v", exitCode)
	os.Exit(exitCode)
}
