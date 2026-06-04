// internal/server/rag.go
package server

import (
	"encoding/json"
	"net/http"

	"auravector/internal/engine"
	"auravector/internal/pipeline"
)

type RagIngestPayload struct {
	Text string `json:"text"`
}

type RagQueryPayload struct {
	Query   string            `json:"query"`
	K       int               `json:"k"`
	Filters map[string]string `json:"filters"` // NEW: Accept JSON filters
}

// handleRagIngest takes a raw document, chunks it, embeds each chunk via Ollama, 
// and safely stores the vector + text in AuraVector-Go.
func (s *Server) handleRagIngest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload RagIngestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Chunk the document using the Sliding Window Chunker (500 chars, 100 overlap)
		chunks := pipeline.ChunkText(payload.Text, 500, 100)
		var ingestedCount int

		for _, chunk := range chunks {
			// 1. Get the embedding from local Ollama
			vector, err := s.embedder.EmbedText(chunk)
			if err != nil {
				continue // In production, we'd log this error
			}

			// 2. Generate a sequential, thread-safe unique ID
			id := s.totalCount.Add(1)

			// 3. Insert the float vector into the high-performance math engine
			s.db.Insert(engine.Vector{ID: id, Values: vector})

			// 4. Insert the raw text chunk into the Document Store
			s.docStore.Insert(id, chunk)
			ingestedCount++
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":          "Document successfully ingested into AuraVector-Go",
			"chunks_processed": ingestedCount,
		})
	}
}

// handleRagQuery embeds the user's question, searches the database for top matches,
// fetches their original text, and uses the LLM to generate a strictly contextual answer.
func (s *Server) handleRagQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload RagQueryPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 1. Convert the user's natural language question into a vector
		queryVector, err := s.embedder.EmbedText(payload.Query)
		if err != nil {
			http.Error(w, "Failed to embed query: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Search AuraVector-Go for the Top-K closest vector matches
		candidates := s.db.Search(queryVector, payload.Filters)

		// NEW FIX: Enforce the Top-K boundary limit requested by the frontend
		// This prevents the engine from flooding the LLM with every chunk in the spatial bucket
		if payload.K > 0 && len(candidates) > payload.K {
			candidates = candidates[:payload.K] 
		}

		// 3. Retrieve the raw text context for those matching IDs
		var contextChunks []string
		for _, cand := range candidates {
			if text, exists := s.docStore.Get(cand); exists {
				contextChunks = append(contextChunks, text)
			}
		}

		// 4. Build the strict RAG prompt
		prompt := pipeline.BuildPrompt(contextChunks, payload.Query)

		// 5. Generate the final answer using the LLM
		answer, err := s.generator.GenerateAnswer(prompt)
		if err != nil {
			http.Error(w, "Failed to generate answer: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 6. Return the intelligence payload to the client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"answer":           answer,
			"contexts_used":    len(contextChunks),
			"search_pool_size": s.totalCount.Load(),
			"contexts":         contextChunks, // NEW: Send the raw text to the UI
		})
	}
}