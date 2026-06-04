// internal/server/router_test.go
package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auravector/internal/engine"
)

func TestInsertAndTelemetry(t *testing.T) {
	// Initialize an empty database and attach it to the server
	db := engine.NewDB(128, 8, 42)
	srv := NewServer(db)

	// 1. Test POST /insert
	insertBody := []byte(`{"id": 1, "values": [0.1, 0.2, 0.3]}`)
	req := httptest.NewRequest(http.MethodPost, "/insert", bytes.NewBuffer(insertBody))
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GET /telemetry to ensure the metric incremented
	reqMetrics := httptest.NewRequest(http.MethodGet, "/telemetry", nil)
	rrMetrics := httptest.NewRecorder()

	srv.ServeHTTP(rrMetrics, reqMetrics)

	var telemetry map[string]interface{}
	json.NewDecoder(rrMetrics.Body).Decode(&telemetry)

	// indexed_vectors is parsed as a float64 from the JSON unmarshaler
	if count, ok := telemetry["indexed_vectors"].(float64); !ok || count != 1 {
		t.Errorf("Expected telemetry to show 1 indexed vector, got %v", telemetry["indexed_vectors"])
	}
}