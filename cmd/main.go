package main

import (
	"bennu.cl/identifier-producer/config"
	"bennu.cl/identifier-producer/internal/api"
	"bennu.cl/identifier-producer/internal/server"
	"bennu.cl/identifier-producer/pkg/kafka"
	"bennu.cl/identifier-producer/version"
	"flag"
	"fmt"
	"k8s.io/klog"
	"math/rand"
	"os"
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

	if *startServer {
		c, err := config.ParseConfig()
		if err != nil {
			klog.Fatalf("Config initialization failed: %s", err)
		}

		idp, err := api.NewIdentifierProducer(c)
		if err != nil {
			klog.Fatalf("Producer initialization failed: %s", err)
		}

		healthz, err := kafka.NewHealthz(c)
		if err != nil {
			klog.Fatalf("Health check initialization failed: %s", err)
		}

		app := server.NewServer(idp, healthz)

		app.Run()
	} else {
		flag.Usage()
	}
}
