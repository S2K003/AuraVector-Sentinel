// cmd/sentinel/main.go
package main

import (
	"fmt"
	"net/http"
	
	// Crucial: This explicitly matches our module name
	"aurasentinel/internal/telemetry"
)

func main() {
	fmt.Println("[*] Booting Aura-Sentinel Reverse Proxy & Threat Hunter...")

	// 1. Create a high-throughput Go channel for our logs with a buffer of 1,000 requests
	telemetryChan := make(chan telemetry.EnterpriseLog, 1000)

	// 2. Start a background Goroutine to listen to the channel (AI Pipeline Placeholder)
	go func() {
		for logEntry := range telemetryChan {
			fmt.Println("\n[+] Real-Time Traffic Intercepted!")
			fmt.Println("    Semantic Context:", logEntry.FlattenForAI())
		}
	}()

	// 3. Create a dummy backend server (Simulating a real web app)
	backendApp := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/admin" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized Access"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the secure web server!"))
	})

	// 4. Wrap the backend with our Zero-Day Interceptor
	secureProxy := telemetry.ThreatInterceptor(backendApp, telemetryChan)

	// 5. Start listening for live web traffic on port 8080
	fmt.Println("[*] Sentinel WAF listening on port 8080...")
	if err := http.ListenAndServe(":8080", secureProxy); err != nil {
		fmt.Printf("[!] Proxy crashed: %v\n", err)
	}
}