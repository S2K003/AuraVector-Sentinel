// internal/telemetry/proxy.go
package telemetry

import (
	"net/http"
	"time"
)

// responseRecorder intercepts the HTTP response to capture the status code and payload size.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytesSent  int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.bytesSent += size
	return size, err
}

// ThreatInterceptor is an HTTP middleware that extracts telemetry in real-time.
func ThreatInterceptor(next http.Handler, telemetryChan chan<- EnterpriseLog) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// 1. Wrap the response writer to spy on the status code and bytes
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// 2. Serve the actual HTTP request
		next.ServeHTTP(rec, req)

		// 3. Calculate response time
		duration := time.Since(start).Milliseconds()

		// 4. Generate the telemetry entry
		logEntry := EnterpriseLog{
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			SourceIP:     req.RemoteAddr,
			Method:       req.Method,
			Path:         req.URL.Path,
			StatusCode:   rec.statusCode,
			BytesSent:    rec.bytesSent,
			UserAgent:    req.UserAgent(),
			ResponseTime: int(duration),
		}

		// 5. Push to the AI engine asynchronously
		select {
		case telemetryChan <- logEntry:
			// Successfully handed off to the AI pipeline
		default:
			// If the channel is full (AI is lagging), drop the telemetry to ensure the web server never crashes.
		}
	})
}