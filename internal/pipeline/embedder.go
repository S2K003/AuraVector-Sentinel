// internal/pipeline/embedder.go
package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmbedRequest represents the JSON payload sent to Ollama.
type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbedResponse represents the JSON payload received from Ollama.
type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// Embedder handles communication with the local Ollama API.
type Embedder struct {
	Host  string
	Model string
}

// NewEmbedder initializes a new Ollama API client.
func NewEmbedder(host, model string) *Embedder {
	if host == "" {
		host = "http://localhost:11434"
	}
	if model == "" {
		model = "mxbai-embed-large"
	}
	return &Embedder{
		Host:  host,
		Model: model,
	}
}

// EmbedText sends a string to Ollama and returns the float32 vector embedding.
func (e *Embedder) EmbedText(text string) ([]float32, error) {
	reqBody := EmbedRequest{
		Model:  e.Model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/embeddings", e.Host)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to reach Ollama (is it running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}

	if len(embedResp.Embedding) == 0 {
		return nil, fmt.Errorf("received empty embedding array from ollama")
	}

	return embedResp.Embedding, nil
}