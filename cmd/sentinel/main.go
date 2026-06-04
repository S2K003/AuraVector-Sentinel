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
	fmt.Println("[*] Booting Aura-Sentinel Edge-WAF & SOC Telemetry...")

	normalLog := telemetry.EnterpriseLog{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		SourceIP:  "127.0.0.1",
		Method:    "GET",
		Path:      "/",
	}
	baselineVector, _ := engine.GenerateSignature(normalLog.FlattenForAI())
	const AnomalyThreshold = 12.0 

	// 1. Initialize our UI Manager & Edge Blocklist
	socManager := telemetry.NewSOCManager()

	// 2. Start the Telemetry Stream for Next.js on port 8081
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/soc/stream", socManager)
		fmt.Println("[*] Live SOC Telemetry Stream running on port 8081")
		http.ListenAndServe(":8081", mux)
	}()

	// 3. The Active Threat Evaluator
	evaluator := func(log telemetry.EnterpriseLog) (bool, string) {
		
		// ==========================================
		// FAST PATH: Check Edge Blocklist O(1)
		// ==========================================
		if socManager.IsBlocked(log.SourceIP) {
			fmt.Printf("[X] Blocked IP %s at the Edge (Bypassed AI)\n", log.SourceIP)
			// Returning true with an empty string drops the connection instantly
			return true, "403 - Connection Terminated by Aura-Sentinel Edge"
		}

		// ==========================================
		// SLOW PATH: AI Vector Math & Trap
		// ==========================================
		semanticContext := log.FlattenForAI()
		vector, err := engine.GenerateSignature(semanticContext)
		if err != nil { return false, "" }

		distance, _ := engine.EuclideanDistance(baselineVector, vector)

		if distance > AnomalyThreshold {
			fmt.Printf("\n[!] ZERO-DAY TRAPPED! Path: %s\n", log.Path)
			
			// 1. Block the IP instantly at the Edge for future requests
			socManager.BlockIP(log.SourceIP)
			
			// 2. Asynchronously generate AI report & push to Next.js UI
			go func() {
				report, _ := engine.GenerateIncidentReport(semanticContext, distance)
				socManager.Notifier <- telemetry.Alert{
					IP:       log.SourceIP,
					Path:     log.Path,
					Distance: distance,
					Report:   report,
				}
			}()
			
			// 3. Tarpit Response (Waste their time once, then never again)
			fakeResponse, _ := engine.GenerateHoneypotResponse(semanticContext, log.Path)
			time.Sleep(3 * time.Second)
			return true, fakeResponse
		}
		return false, ""
	}

	// 4. Wrap App and Start WAF on 8080
	backendApp := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the secure web server!"))
	})
	secureProxy := telemetry.ActiveThreatInterceptor(backendApp, evaluator)

	fmt.Println("[*] Sentinel WAF listening on port 8080...")
	http.ListenAndServe(":8080", secureProxy)
}