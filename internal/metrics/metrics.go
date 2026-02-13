package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		},
		[]string{"method", "endpoint"},
	)

	modelPredictionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "model_predictions_total",
			Help: "Total number of model predictions",
		},
		[]string{"model_id", "status"},
	)

	modelInferenceDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "model_inference_duration_seconds",
			Help:    "Model inference duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0},
		},
		[]string{"model_id"},
	)

	modelsLoaded = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "models_loaded",
			Help: "Number of models currently loaded",
		},
	)
)

func RecordHTTPRequest(method, endpoint string, status int, duration float64) {
	httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func RecordPrediction(modelID string, success bool, duration float64) {
	if success {
		modelPredictionsTotal.WithLabelValues(modelID, "success").Inc()
	} else {
		modelPredictionsTotal.WithLabelValues(modelID, "error").Inc()
	}

	modelInferenceDuration.WithLabelValues(modelID).Observe(duration)
}

func SetModelsLoaded(count int) {
	modelsLoaded.Set(float64(count))
}
