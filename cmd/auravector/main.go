// cmd/auravector/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // <-- Added the missing os import

	"auravector/internal/engine"
	"auravector/internal/pipeline"
	"auravector/internal/server"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Booting AuraVector-Go RAG Intelligence Engine...")

	// 1. Ensure the persistent data directory exists inside the container
	os.MkdirAll("/app/data", 0755)

	// 2. Initialize Core DB and execute Vector State Recovery
	db := engine.NewDB(1024, 2, 42)
	if err := db.Recover("/app/data/vectors.vdb"); err != nil {
		log.Fatalf("Failed to recover vector database: %v", err)
	}

	// 3. Initialize Persistent Document Store (Recovers automatically on boot)
	docStore := engine.NewDocumentStore("/app/data/docs.jsonl")

	// 4. Initialize AI Pipeline
	embedder := pipeline.NewEmbedder("http://host.docker.internal:11434", "")
	generator := pipeline.NewGenerator("http://host.docker.internal:11434", "")

	// 5. Wrap the execution components in the Server multiplexer
	srv := server.NewServer(db, docStore, embedder, generator)

	// 6. Use our new public setter to sync the server ID counter
	srv.SetTotalCount(docStore.MaxID)

	port := ":8080"
	fmt.Printf("RAG Engine live and listening on http://localhost%s\n", port)

	// Start the server
	log.Fatal(http.ListenAndServe(port, corsMiddleware(srv)))
}