// internal/pipeline/generator_test.go
package pipeline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	chunks := []string{"AuraVector-Go is incredibly fast.", "It uses zero-allocation memory mapping."}
	query := "How does it manage memory?"
	
	prompt := BuildPrompt(chunks, query)

	// Verify structural integrity of the prompt
	if !strings.Contains(prompt, "--- Chunk 1 ---") || !strings.Contains(prompt, "AuraVector-Go is incredibly fast.") {
		t.Errorf("Prompt is missing the first context chunk.")
	}
	if !strings.Contains(prompt, "--- Chunk 2 ---") || !strings.Contains(prompt, "zero-allocation") {
		t.Errorf("Prompt is missing the second context chunk.")
	}
	if !strings.Contains(prompt, query) {
		t.Errorf("Prompt is missing the user's query.")
	}
}

func TestGenerateAnswer(t *testing.T) {
	// Mock the Ollama generation server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}

		mockResponse := GenerateResponse{
			Response: "This is a strictly mocked LLM answer.",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	generator := NewGenerator(mockServer.URL, "dummy-model")
	result, err := generator.GenerateAnswer("Dummy prompt string")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "This is a strictly mocked LLM answer." {
		t.Errorf("Decoding failed. Expected mock answer, got %s", result)
	}
}