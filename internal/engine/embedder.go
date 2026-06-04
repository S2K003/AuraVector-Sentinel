// internal/engine/embedder.go
package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// EmbeddingRequest formats the payload for Ollama.
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse captures the 768-D float array from Ollama.
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerateSignature sends the semantic log to the AI and returns the mathematical vector.
func GenerateSignature(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI request: %v", err)
	}

	// Make a lightning-fast local API call to Ollama
	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not reach Ollama (Is it running?): %v", err)
	}
	defer resp.Body.Close()

	var resBody EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("failed to decode AI response: %v", err)
	}

	if len(resBody.Embedding) == 0 {
		return nil, fmt.Errorf("AI returned an empty mathematical signature")
	}

	return resBody.Embedding, nil
}