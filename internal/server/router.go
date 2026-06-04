// internal/server/router.go
package server

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync/atomic"

	"auravector/internal/engine"
	"auravector/internal/pipeline"
)

// Server wraps the HTTP multiplexer and the core database engine + RAG pipeline.
type Server struct {
	db         *engine.DB
	docStore   *engine.DocumentStore
	embedder   *pipeline.Embedder
	generator  *pipeline.Generator
	mux        *http.ServeMux
	totalCount atomic.Uint64
}

// NewServer initializes the HTTP routing layer with all RAG dependencies.
func NewServer(db *engine.DB, docStore *engine.DocumentStore, embedder *pipeline.Embedder, generator *pipeline.Generator) *Server {
	s := &Server{
		db:        db,
		docStore:  docStore,
		embedder:  embedder,
		generator: generator,
		mux:       http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
    // Existing routes...
	s.mux.HandleFunc("/insert", s.handleInsert())
	s.mux.HandleFunc("/insert/batch", s.handleInsertBatch())
	s.mux.HandleFunc("/search", s.handleSearch())
	s.mux.HandleFunc("/telemetry", s.handleTelemetry())
	
	// New RAG Routes
	s.mux.HandleFunc("/rag/ingest", s.handleRagIngest())
	s.mux.HandleFunc("/rag/query", s.handleRagQuery())
}

// --- Payload Structures ---

type InsertPayload struct {
	ID     uint64    `json:"id"`
	Values []float32 `json:"values"`
}

type SearchPayload struct {
	Vector []float32 `json:"vector"`
	K      int       `json:"k"`
}

// --- Handlers ---


func (s *Server) SetTotalCount(count uint64) {
	s.totalCount.Store(count)
}

func (s *Server) handleInsert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload InsertPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.db.Insert(engine.Vector{ID: payload.ID, Values: payload.Values})
		s.totalCount.Add(1)
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) handleInsertBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payloads []InsertPayload
		if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, p := range payloads {
			s.db.Insert(engine.Vector{ID: p.ID, Values: p.Values})
		}
		
		s.totalCount.Add(uint64(len(payloads)))
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload SearchPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Execute the ANN search using the specific bounds
		candidates := s.db.Search(payload.Vector, nil) // Pass nil to ignore filters

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"candidates": candidates,
			"count":      len(candidates),
		})
	}
}

func (s *Server) handleTelemetry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		telemetry := map[string]interface{}{
			"indexed_vectors": s.totalCount.Load(),
			"alloc_bytes":     m.Alloc,
			"sys_bytes":       m.Sys,
			"num_gc":          m.NumGC,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(telemetry)
	}
}