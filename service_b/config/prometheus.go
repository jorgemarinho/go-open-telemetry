package config

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheusMetrics(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(port, nil)
}
