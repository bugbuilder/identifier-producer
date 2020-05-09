package handlers

import (
	"bennu.cl/identifier-producer/pkg/api"
	"bennu.cl/identifier-producer/pkg/kafka"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"syscall"
)

func Producer(ids api.IdentifierService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id api.Identifier

		if err := c.BindJSON(&id); err == nil {
			key, err := ids.Save(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Producer failed",
					"error":   err.Error(),
				})

				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			} else {
				c.JSON(http.StatusCreated, gin.H{"key": key})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error:": "BadRequest",
				"mesage": "The request could not be understood by the server due to malformed syntax",
			})
		}
	}
}

func Healthz(h kafka.Healthz) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.AvailableCluster(); err != nil {
			klog.Errorf("%s", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "AvailableCluster failed",
				"error":   err.Error(),
			})
		}

		if err := h.AvailablePartitions(); err != nil {
			klog.Errorf("%s", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "AvailablePartitions failed",
				"error":   err.Error(),
			})
		}
	}
}
