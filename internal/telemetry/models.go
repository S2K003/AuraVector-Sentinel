// internal/telemetry/models.go
package telemetry

import "fmt"

// EnterpriseLog represents the structured metadata of an HTTP request.
type EnterpriseLog struct {
	Timestamp    string `json:"timestamp"`
	SourceIP     string `json:"source_ip"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	StatusCode   int    `json:"status_code"`
	BytesSent    int    `json:"bytes_sent"`
	UserAgent    string `json:"user_agent"`
	ResponseTime int    `json:"response_time_ms"`
}

// FlattenForAI converts the structured JSON into a contextual semantic sentence.
func (l EnterpriseLog) FlattenForAI() string {
	return fmt.Sprintf("Traffic from IP %s executed %s on %s. Status: %d. Transferred %d bytes in %d ms. User-Agent: %s",
		l.SourceIP, l.Method, l.Path, l.StatusCode, l.BytesSent, l.ResponseTime, l.UserAgent)
}