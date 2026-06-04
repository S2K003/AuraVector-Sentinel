// internal/telemetry/soc.go
package telemetry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Alert represents the data sent to the Next.js UI
type Alert struct {
	IP       string  `json:"ip"`
	Path     string  `json:"path"`
	Distance float64 `json:"distance"`
	Report   string  `json:"report"`
}

// SOCManager handles the Live UI stream and the Adaptive Blocklist
type SOCManager struct {
	BlockedIPs sync.Map // Thread-safe fast-path cache
	Notifier   chan Alert
	clients    map[chan Alert]bool
	mu         sync.Mutex
}

func NewSOCManager() *SOCManager {
	soc := &SOCManager{
		Notifier: make(chan Alert, 100),
		clients:  make(map[chan Alert]bool),
	}
	go soc.listen()
	return soc
}

// BlockIP caches the hacker's IP in O(1) time
func (s *SOCManager) BlockIP(ip string) {
	s.BlockedIPs.Store(ip, true)
}

// IsBlocked checks the fast-path cache instantly
func (s *SOCManager) IsBlocked(ip string) bool {
	_, exists := s.BlockedIPs.Load(ip)
	return exists
}

// listen broadcasts alerts to all connected Next.js clients
func (s *SOCManager) listen() {
	for event := range s.Notifier {
		s.mu.Lock()
		for clientChan := range s.clients {
			clientChan <- event
		}
		s.mu.Unlock()
	}
}

// ServeHTTP handles the SSE connection and CORS for Next.js
func (s *SOCManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// CORS Headers are required because Next.js runs on a different port (3000)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan Alert)
	
	s.mu.Lock()
	s.clients[messageChan] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, messageChan)
		s.mu.Unlock()
	}()

	flusher, _ := w.(http.Flusher)

	for {
		select {
		case alert := <-messageChan:
			data, _ := json.Marshal(alert)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush() // Push to UI instantly
		case <-req.Context().Done():
			return
		}
	}
}