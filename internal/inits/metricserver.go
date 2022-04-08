package inits

import (
	"github.com/jaczerob/madamchuckle/internal/services/metrics"

	log "github.com/sirupsen/logrus"
)

func InitMetricsServer() (m *metrics.MetricsServer, err error) {
	m, err = metrics.NewMetricsServer()
	if err != nil {
		return
	}

	go func() {
		if err = m.Start(); err != nil {
			log.WithError(err).Fatal(err)
		}
	}()

	return
}
