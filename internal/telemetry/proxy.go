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

// ThreatEvaluator is a callback function that evaluates a log and returns (isThreat, fakeResponse)
type ThreatEvaluator func(log EnterpriseLog) (bool, string)

// ActiveThreatInterceptor evaluates threats BEFORE they hit the real web server.
func ActiveThreatInterceptor(next http.Handler, evaluate ThreatEvaluator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Capture inbound metadata (Status/Bytes are 0 because the response hasn't happened yet)
		logEntry := EnterpriseLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			SourceIP:  req.RemoteAddr,
			Method:    req.Method,
			Path:      req.URL.Path,
			UserAgent: req.UserAgent(),
		}

		// 1. Evaluate the threat mathematically in real-time
		isThreat, fakeResponse := evaluate(logEntry)

		if isThreat {
			// THE TARPIT: Feed them the fake AI response, masquerading as a success
			w.WriteHeader(http.StatusOK) // Return a 200 OK to trick automated scanners!
			w.Write([]byte(fakeResponse))
			return
		}

		// 2. Normal traffic flows through to the real backend safely
		next.ServeHTTP(w, req)
	})
}