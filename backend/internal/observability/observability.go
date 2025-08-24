package observability

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.Tracer("commerce/backend")

var (
	promRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_client_requests_total",
		Help: "Total outbound HTTP client requests",
	}, []string{"service", "method", "status"})

	promLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_client_latency_seconds",
		Help:    "Outbound HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method"})
)

// EnsureRequestID makes sure a request has an X-Request-Id header.
func EnsureRequestID(req *http.Request) {
	if req.Header.Get("X-Request-Id") == "" {
		req.Header.Set("X-Request-Id", uuid.NewString())
	}
}

// InjectTraceHeaders injects the current trace context into the outgoing HTTP headers.
func InjectTraceHeaders(ctx context.Context, req *http.Request) {
	if ctx == nil {
		ctx = context.TODO()
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}

// StartSpan starts a new span for an outbound HTTP call and returns a context and a finish function.
// The finish function should be called with the HTTP status and error (if any).
func StartSpan(ctx context.Context, spanName, service, method, url string) (context.Context, func(status int, err error)) {
	ctx, span := tracer.Start(ctx, spanName, trace.WithAttributes(
		attribute.String("service", service),
		attribute.String("http.method", method),
		attribute.String("http.url", url),
	))
	start := time.Now()

	return ctx, func(status int, err error) {
		dur := time.Since(start)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.SetAttributes(attribute.Int("http.status_code", status))

		// emit Prometheus metrics and fallback log
		promRequests.WithLabelValues(service, method, fmt.Sprintf("%d", status)).Inc()
		promLatency.WithLabelValues(service, method).Observe(dur.Seconds())
		RecordServiceRequest(service, method, url, status, dur)
		span.End()
	}
}

// RecordServiceRequest logs a simple trace for outgoing service requests.
func RecordServiceRequest(service, method, url string, status int, duration time.Duration) {
	log.Printf("observability: service=%s method=%s url=%s status=%d duration=%s", service, method, url, status, duration)
}
