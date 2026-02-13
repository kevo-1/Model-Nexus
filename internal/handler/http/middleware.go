package http

import (
    "net/http"
    "time"
    
    "github.com/kevo-1/model-serving-platform/internal/metrics"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
            next.ServeHTTP(w, r)
            return
        }
		
        start := time.Now()
        
        wrapped := newResponseWriter(w)
        
        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()
        metrics.RecordHTTPRequest(r.Method, r.URL.Path, wrapped.statusCode, duration)
    })
}