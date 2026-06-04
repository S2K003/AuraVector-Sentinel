// cmd/sentinel/main.go
package main

import (
	"fmt"
	"net/http"

	"aurasentinel/internal/engine"
	"aurasentinel/internal/telemetry"
)

func main() {
	fmt.Println("[*] Booting Aura-Sentinel Reverse Proxy & Threat Hunter...")

	telemetryChan := make(chan telemetry.EnterpriseLog, 1000)

	// Background AI Pipeline
	go func() {
		for logEntry := range telemetryChan {
			fmt.Println("\n[+] Real-Time Traffic Intercepted!")
			semanticContext := logEntry.FlattenForAI()
			fmt.Println("    Context:", semanticContext)
			
			// 1. Send the context to the AI Embedder
			fmt.Println("    [AI] Generating 768-Dimensional Threat Signature...")
			vector, err := engine.GenerateSignature(semanticContext)
			if err != nil {
				fmt.Printf("    [!] Pipeline Error: %v\n", err)
				continue
			}

			// 2. Prove the math works by printing the first 3 dimensions
			fmt.Printf("    [+] Signature Established: [%.4f, %.4f, %.4f ...]\n", vector[0], vector[1], vector[2])
		}
	}()

	backendApp := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/admin" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized Access"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the secure web server!"))
	})

	secureProxy := telemetry.ThreatInterceptor(backendApp, telemetryChan)

	fmt.Println("[*] Sentinel WAF listening on port 8080...")
	if err := http.ListenAndServe(":8080", secureProxy); err != nil {
		fmt.Printf("[!] Proxy crashed: %v\n", err)
	}
}