package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标集合
type Metrics struct {
	HTTPRequestsTotal      prometheus.CounterVec
	HTTPRequestDuration    prometheus.HistogramVec
	HTTPInflightRequests   prometheus.GaugeVec
	InferenceRequestsTotal prometheus.CounterVec
	InferenceLatency       prometheus.HistogramVec
	InferenceTTFT          prometheus.HistogramVec
	PromptTokensTotal      prometheus.CounterVec
	CompletionTokensTotal  prometheus.CounterVec
	BackendRequestsTotal   prometheus.CounterVec
	BackendErrorsTotal     prometheus.CounterVec
	AuthFailuresTotal      prometheus.CounterVec
	RateLimitRejected      prometheus.CounterVec
	SessionsActive         prometheus.Gauge
	SessionsEvicted        prometheus.CounterVec
}

// NewMetrics 创建指标集合
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_http_requests_total",
				Help: "Total HTTP requests",
			},
			[]string{"route", "method", "status"},
		),
		HTTPRequestDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "llamago_http_request_duration_seconds",
				Help:    "HTTP request duration",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"route", "method", "status"},
		),
		HTTPInflightRequests: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "llamago_http_inflight_requests",
				Help: "Inflight HTTP requests",
			},
			[]string{"route"},
		),
		InferenceRequestsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_inference_requests_total",
				Help: "Total inference requests",
			},
			[]string{"model", "backend", "stream"},
		),
		InferenceLatency: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "llamago_inference_latency_seconds",
				Help:    "Inference latency",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"model", "backend", "stream"},
		),
		InferenceTTFT: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "llamago_inference_ttft_seconds",
				Help:    "Time to first token",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"model", "backend"},
		),
		PromptTokensTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_inference_prompt_tokens_total",
				Help: "Total prompt tokens",
			},
			[]string{"model", "backend"},
		),
		CompletionTokensTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_inference_completion_tokens_total",
				Help: "Total completion tokens",
			},
			[]string{"model", "backend"},
		),
		BackendRequestsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_backend_requests_total",
				Help: "Total backend requests",
			},
			[]string{"backend"},
		),
		BackendErrorsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_backend_errors_total",
				Help: "Total backend errors",
			},
			[]string{"backend", "error_type"},
		),
		AuthFailuresTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_auth_failures_total",
				Help: "Total auth failures",
			},
			[]string{"reason"},
		),
		RateLimitRejected: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_rate_limit_rejected_total",
				Help: "Total rate limit rejections",
			},
			[]string{"scope"},
		),
		SessionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "llamago_sessions_active",
				Help: "Active sessions",
			},
		),
		SessionsEvicted: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llamago_sessions_evicted_total",
				Help: "Total evicted sessions",
			},
			[]string{"reason"},
		),
	}
}
