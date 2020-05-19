package handlers

import (
	"bennu.cl/identifier-producer/pkg/kafka"
	"k8s.io/klog"
	"net/http"
)

// todo: move to HealthCheckServer
func Healthz(h kafka.Kafka) http.HandlerFunc {
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
