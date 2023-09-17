package controller

import (
	"github.com/bhatti/PlexAuthZ/internal/web"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// MetricsController for displaying metrics.
type MetricsController struct {
}

// NewMetricsController instantiates controller for displaying metrics.
func NewMetricsController(
	webserver web.Server) (*MetricsController, error) {
	webserver.GET("/metrics", web.WrapHandler(promhttp.Handler()))
	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		logrus.WithFields(logrus.Fields{
			"Component": "MetricsController",
			"Error":     err,
		}).Error("failed to register prometheus collector")
	}
	return &MetricsController{}, nil
}
