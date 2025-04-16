package urls

import (
	"fmt"

	"github.com/mohammadne/fesghel/internal/entities"
	metrics_pkg "github.com/mohammadne/fesghel/pkg/observability/metrics"
)

type metrics struct {
	Counter   metrics_pkg.Counter
	Histogram metrics_pkg.Histogram
}

func newMetrics() (metrics *metrics, err error) {
	var prefix = "urls"

	counterName := prefix + "_counter"
	counterLabels := []string{"method", "status"}
	metrics.Counter, err = metrics_pkg.RegisterCounter(counterName, entities.Namespace, entities.System, counterLabels)
	if err != nil {
		return nil, fmt.Errorf("error while registering counter vector: %v", err)
	}

	histogramName := prefix + "_histogram"
	histogramLabels := []string{"method"}
	metrics.Histogram, err = metrics_pkg.RegisterHistogram(histogramName, entities.Namespace, entities.System, histogramLabels)
	if err != nil {
		return nil, fmt.Errorf("error while registering histogram vector: %v", err)
	}

	return metrics, nil
}
