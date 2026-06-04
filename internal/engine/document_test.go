// internal/engine/document_test.go
package engine

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestDocumentStoreConcurrent(t *testing.T) {
	// Use Go's built-in TempDir which automatically cleans up after the test finishes
	tempFile := filepath.Join(t.TempDir(), "test_docs.jsonl")
	
	// Pass the temporary file path to our upgraded DocumentStore
	store := NewDocumentStore(tempFile)
	var wg sync.WaitGroup

	// Concurrently insert 1000 documents
	for i := uint64(0); i < 1000; i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			store.Insert(id, "mock text chunk")
		}(i)
	}

	// Concurrently read them back
	for i := uint64(0); i < 1000; i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			_, _ = store.Get(id)
		}(i)
	}

	wg.Wait()

	// Verify an expected chunk exists
	if text, exists := store.Get(500); !exists || text != "mock text chunk" {
		t.Fatalf("Failed to retrieve expected document text.")
	}
}