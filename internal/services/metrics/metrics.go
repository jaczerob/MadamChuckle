package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	PopulationTracker = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "madamchuckle_population_tracker",
		Help: "The population over time of Toontown Rewritten",
	})

	PopulationByHourTracker = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "madamchuckle_populationbyhour_tracker",
		Help:    "The population average for each hour of Toontown Rewritten",
		Buckets: prometheus.LinearBuckets(100, 100, 60),
	}, []string{"hour"})
)

type MetricsServer struct {
	server *http.Server
}

func NewMetricsServer() (m *MetricsServer, err error) {
	m = new(MetricsServer)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	m.server = &http.Server{
		Addr:    viper.GetString("prometheus.addr"),
		Handler: mux,
	}

	return
}

func (m *MetricsServer) Start() error {
	return m.server.ListenAndServe()
}

func (m *MetricsServer) Shutdown() error {
	return m.server.Shutdown(context.Background())
}
