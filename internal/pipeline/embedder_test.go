// internal/pipeline/embedder_test.go
package pipeline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbedText(t *testing.T) {
	// 1. Create a mock HTTP server that simulates Ollama's exact response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method and URI
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("Expected path /api/embeddings, got %s", r.URL.Path)
		}

		// Send back a dummy float32 array
		mockResponse := EmbedResponse{
			Embedding: []float32{0.1, 0.2, 0.3, 0.4, 0.5},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	// 2. Point our Embedder to the mock server URL instead of localhost:11434
	embedder := NewEmbedder(mockServer.URL, "dummy-model")

	// 3. Execute the function
	result, err := embedder.EmbedText("Test chunk string")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 5 {
		t.Fatalf("Expected embedding length 5, got %d", len(result))
	}

	if result[0] != 0.1 || result[4] != 0.5 {
		t.Errorf("Embedding decoding failed. Got %v", result)
	}
}