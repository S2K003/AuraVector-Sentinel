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
	fmt.Println("[*] Booting Aura-Sentinel AI Tarpit & Active WAF...")

	fmt.Println("[*] Training AI Baseline Profile...")
	normalLog := telemetry.EnterpriseLog{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		SourceIP:  "127.0.0.1",
		Method:    "GET",
		Path:      "/",
		UserAgent: "Mozilla/5.0",
	}
	baselineVector, _ := engine.GenerateSignature(normalLog.FlattenForAI())
	fmt.Println("[+] Baseline Locked. System Armed.")
	
	const AnomalyThreshold = 12.0 

	// 1. The Active Threat Evaluator (Inline Trap)
	evaluator := func(log telemetry.EnterpriseLog) (bool, string) {
		semanticContext := log.FlattenForAI()
		
		// Convert inbound request to Math
		vector, err := engine.GenerateSignature(semanticContext)
		if err != nil {
			return false, ""
		}

		// Calculate Distance
		distance, _ := engine.EuclideanDistance(baselineVector, vector)

		// Trigger the Tarpit
		if distance > AnomalyThreshold {
			fmt.Printf("\n[!] ZERO-DAY TRAPPED! (Distance: %.2f) Path: %s\n", distance, log.Path)
			fmt.Println("    Action: Rerouting Hacker to GenAI Tarpit...")
			
			// Generate Fake Environment
			fakeResponse, err := engine.GenerateHoneypotResponse(semanticContext, log.Path)
			if err != nil {
				return true, "{\"error\": \"database timeout\"}"
			}
			
			fmt.Println("    [Tarpit] Generating fake vulnerability payload...")
			fmt.Println("    [Tarpit] Delaying response by 3 seconds to drain attacker threads...")
			
			// Tarpit delay: Waste the hacker's connection threads
			time.Sleep(3 * time.Second)
			return true, fakeResponse
		}

		fmt.Printf("[OK] Traffic Normal. (Distance: %.2f)\n", distance)
		return false, ""
	}

	// 2. The Real Web Application
	backendApp := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the secure web server! Your real data is here."))
	})

	// 3. Wrap the App with the Active Interceptor
	secureProxy := telemetry.ActiveThreatInterceptor(backendApp, evaluator)

	fmt.Println("[*] Sentinel Tarpit listening on port 8080...")
	http.ListenAndServe(":8080", secureProxy)
}