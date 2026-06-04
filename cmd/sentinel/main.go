// cmd/sentinel/main.go
package main

import (
	"fmt"
	"net/http"
	"time"

	"aurasentinel/internal/engine"
	"aurasentinel/internal/telemetry"
)

func main() {
	fmt.Println("[*] Booting Aura-Sentinel Zero-Day Threat Hunter...")

	// 1. Establish the Baseline "Normal" Traffic Signature
	// In a production system, this would be the average of 10,000 normal logs.
	fmt.Println("[*] Training AI Baseline Profile...")
	normalLog := telemetry.EnterpriseLog{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		SourceIP:     "127.0.0.1",
		Method:       "GET",
		Path:         "/",
		StatusCode:   200,
		BytesSent:    500,
		UserAgent:    "Mozilla/5.0",
		ResponseTime: 45,
	}
	baselineVector, err := engine.GenerateSignature(normalLog.FlattenForAI())
	if err != nil {
		fmt.Printf("[!] Fatal: Could not establish baseline: %v\n", err)
		return
	}
	fmt.Println("[+] Baseline Locked. System Armed.")
	
	// Define our mathematical boundary. For normalized vectors, > 10.0 usually implies high deviation.
	// You can adjust this based on testing!
	const AnomalyThreshold = 12.0 

	telemetryChan := make(chan telemetry.EnterpriseLog, 1000)

	// Background AI Pipeline (The Trap)
	go func() {
		for logEntry := range telemetryChan {
			fmt.Println("\n-------------------------------------------------")
			fmt.Printf("[+] Intercepted: %s %s from %s\n", logEntry.Method, logEntry.Path, logEntry.SourceIP)
			
			semanticContext := logEntry.FlattenForAI()
			
			// 1. Generate Signature
			vector, err := engine.GenerateSignature(semanticContext)
			if err != nil {
				continue
			}

			// 2. Calculate Distance to Baseline
			distance, _ := engine.EuclideanDistance(baselineVector, vector)

			// 3. Trigger the Trap
			if distance > AnomalyThreshold {
				fmt.Printf("[!] ZERO-DAY ANOMALY DETECTED! (Distance: %.2f)\n", distance)
				fmt.Printf("    Reason: Mathematical deviation exceeds threshold of %.2f\n", AnomalyThreshold)
				fmt.Println("    Action: Initiating Block & Incident Report...")
			} else {
				fmt.Printf("[OK] Traffic Normal. (Distance: %.2f)\n", distance)
			}
			fmt.Println("-------------------------------------------------")
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
	http.ListenAndServe(":8080", secureProxy)
}