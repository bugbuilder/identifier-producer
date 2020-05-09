package server

import (
	"bennu.cl/identifier-producer/internal/handlers"
	"bennu.cl/identifier-producer/pkg/api"
	"bennu.cl/identifier-producer/pkg/kafka"
	"context"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Identifier struct {
	router *gin.Engine

	identifier api.IdentifierService
	healthz    kafka.Healthz
}

func NewServer(ids api.IdentifierService, healthz kafka.Healthz) *Identifier {
	srv := &Identifier{}

	r := gin.New()

	srv.router = r
	srv.identifier = ids
	srv.healthz = healthz
	srv.setRouters()

	return srv
}

func (srv *Identifier) setRouters() {
	srv.router.GET("/healthz", handlers.Healthz(srv.healthz))
	srv.router.POST("/identifier", handlers.Producer(srv.identifier))
}

func (srv *Identifier) Run() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: srv.router,
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
