package sql

import (
	"github.com/najeira/goutils/metrics"
)

var (
	Metrics *metrics.MetricsDB
)

func init() {
	Metrics = metrics.NewMetricsDB()
}
