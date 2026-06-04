// internal/engine/analyzer.go
package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GenerateRequest formats the prompt for the generative LLM.
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse captures the AI's generated text.
type GenerateResponse struct {
	Response string `json:"response"`
}

// GenerateIncidentReport acts as an automated SOC Analyst.
func GenerateIncidentReport(threatContext string, distance float64) (string, error) {
	// We inject the anomalous log into a strict system prompt
	systemPrompt := fmt.Sprintf(`You are an elite Cybersecurity Analyst AI. 
A Zero-Day anomaly was just trapped by our vector firewall with a mathematical severity score of %.2f.
Analyze this raw intercepted web request: "%s"

Provide a highly concise, 3-bullet-point Incident Report:
1. Suspected Attack Type (e.g., Brute Force, SQLi, Reconnaissance)
2. Hacker Intent
3. Recommended Firewall Rule to prevent future attempts.`, distance, threatContext)

	reqBody := GenerateRequest{
		Model:  "llama3.2", // We use a fast, generative LLM here
		Prompt: systemPrompt,
		Stream: false,
	}
	
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("could not reach AI Analyst: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var resBody GenerateResponse
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		return "", fmt.Errorf("failed to decode AI report")
	}

	return resBody.Response, nil
}