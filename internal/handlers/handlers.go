package handlers

import (
	"bennu.cl/identifier-producer/pkg/api"
	"bennu.cl/identifier-producer/pkg/kafka"
	"encoding/json"
	"fmt"
	"k8s.io/klog"
	"net/http"
	"syscall"
)

type Message struct {
}

func Producer(ids api.IdentifierService) http.HandlerFunc {
	var id api.Identifier

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(id); err != nil {
				if key, err := ids.Save(id); err == nil {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintln(w, fmt.Sprintf("{\"key\":\"%s\"}", key))
				} else {
					http.Error(w, "Invalid request method", http.StatusInternalServerError)
					// any error we stop the server to stop receiving request
					syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				}
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func Healthz(h kafka.Healthz) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := h.AvailableCluster(); err != nil {
				klog.Errorf("%s", err)
				http.Error(w, "AvailableCluster failed", http.StatusServiceUnavailable)
			}

			if err := h.AvailablePartitions(); err != nil {
				klog.Errorf("%s", err)
				http.Error(w, "AvailablePartitions failed", http.StatusServiceUnavailable)
			}
		}
	}
}
