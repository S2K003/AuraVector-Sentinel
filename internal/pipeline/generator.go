// internal/pipeline/generator.go
package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GenerateRequest represents the JSON payload sent to Ollama for text generation.
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // We set this to false to get a single complete response
}

// GenerateResponse represents the JSON payload received from Ollama.
type GenerateResponse struct {
	Response string `json:"response"`
}

// Generator handles communication with the local Ollama LLM.
type Generator struct {
	Host  string
	Model string
}

// NewGenerator initializes a new Ollama API client for text generation.
func NewGenerator(host, model string) *Generator {
	if host == "" {
		host = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3" // Default to standard Llama 3 model
	}
	return &Generator{
		Host:  host,
		Model: model,
	}
}

// BuildPrompt constructs the strict RAG instruction template.
func BuildPrompt(contextChunks []string, query string) string {
	prompt := "You are a helpful and precise assistant. Use the following retrieved context chunks to answer the user's question. If the answer is not contained in the context, say 'I do not have enough information to answer that based on the provided documents.' Do NOT invent an answer outside of this context.\n\nContext:\n"
	
	for i, chunk := range contextChunks {
		prompt += fmt.Sprintf("--- Chunk %d ---\n%s\n\n", i+1, chunk)
	}
	
	prompt += "Question: " + query + "\nAnswer:"
	return prompt
}

// GenerateAnswer sends the combined prompt string to Ollama and returns the LLM's response.
func (g *Generator) GenerateAnswer(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  g.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", g.Host)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to reach Ollama (is it running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w", err)
	}

	return genResp.Response, nil
}