package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Sessions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "radioserver_sessions",
		Help: "The total number of sessions",
	})
)
