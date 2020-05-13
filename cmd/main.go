package main

import (
	"bennu.cl/identifier-producer/api/server"
	"bennu.cl/identifier-producer/api/server/identifier"
	"bennu.cl/identifier-producer/config"
	identifier2 "bennu.cl/identifier-producer/pkg/identifier"
	"bennu.cl/identifier-producer/version"
	"flag"
	"fmt"
	"k8s.io/klog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	klog.InitFlags(nil)

	rand.Seed(time.Now().UnixNano())

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s -server \n\n", os.Args[0])

		flag.PrintDefaults()
	}
}

func main() {
	defer klog.Flush()

	showVersion := flag.Bool("version", false, "Print the version information")
	startServer := flag.Bool("server", false, "Start server")

	flag.Parse()

	if *showVersion {
		fmt.Println(version.NewInfo().Print())
		os.Exit(0)
	}
	if !*startServer {
		flag.Usage()
		os.Exit(0)
	}

	c, err := config.ParseConfig()
	if err != nil {
		klog.Fatalf("Config initialization failed: %s", err)
	}

	ids, err := identifier2.NewIdentifierProducer(c)
	if err != nil {
		klog.Fatalf("Producer initialization failed: %s", err)
	}

	s := server.NewServer()
	routers := []server.Router{
		identifier.NewRouter(ids),
	}

	s.InitRouter(routers...)

	errRun := make(chan error)
	go s.Run(errRun)

	go gracefulShutdown(s)

	errAPISrv := <-errRun
	if errAPISrv != nil {
		klog.Fatalf("shutting down due to API Server error: %s", errAPISrv)
		os.Exit(1)
	}
}

func gracefulShutdown(s *server.Server) {
	exitCode := 0
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	errShutdown := make(chan error)
	go s.Shutdown(errShutdown)

	errAPISrv := <-errShutdown
	if errAPISrv != nil {
		klog.Fatalf("shutting down due to API Server error: %s", errAPISrv)
		exitCode = 1
	}

	klog.Infof("awaiting kafka shutting down")
	time.Sleep(1 * time.Second)

	klog.Infof("exiting with %v", exitCode)
	os.Exit(exitCode)
}
