// internal/engine/tarpit.go
package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GenerateHoneypotResponse asks LLaMA to generate a fake, believable server response.
func GenerateHoneypotResponse(threatContext, path string) (string, error) {
	// The prompt instructs the AI to behave like a vulnerable, broken server
	systemPrompt := fmt.Sprintf(`You are an elite AI Honeypot (Tarpit) protecting a server. 
A malicious actor just sent this payload to the '%s' endpoint: "%s"

Generate a highly realistic, fake HTTP response body that makes the hacker think their attack was successful or revealed sensitive internal data. 
- If they tried an SQL injection, output a fake SQL syntax error containing fake table names. 
- If they scanned for an admin panel, output a fake, hardcoded HTML login form.
- If they tried directory traversal, output fake internal file paths.

DO NOT include any conversational text or pleasantries. ONLY output the raw fake payload.`, path, threatContext)

	reqBody := GenerateRequest{
		Model:  "llama3.2",
		Prompt: systemPrompt,
		Stream: false,
	}
	
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "{\"error\": \"internal database connection timeout\"}", fmt.Errorf("could not reach AI: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var resBody GenerateResponse
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		return "{\"error\": \"fatal memory exception\"}", fmt.Errorf("failed to decode AI report")
	}

	return resBody.Response, nil
}