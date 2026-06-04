// internal/pipeline/chunker_test.go
package pipeline

import (
	"testing"
)

func TestChunkText(t *testing.T) {
	text := "AuraVector is a fast engine. We are adding RAG to it."
	// Total length: 53 characters
	
	// Chunk size 20, Overlap 5
	// Step size = 15
	chunks := ChunkText(text, 20, 5)

	if len(chunks) == 0 {
		t.Fatalf("Expected chunks, got none")
	}

	// Chunk 1: "AuraVector is a fast" (indices 0-20)
	expectedFirst := "AuraVector is a fast"
	if chunks[0] != expectedFirst {
		t.Errorf("Expected first chunk %q, got %q", expectedFirst, chunks[0])
	}

	// Chunk 2 should overlap by 5 characters. 
	// The last 5 chars of Chunk 1 are "fast" (plus a leading space).
	// So Chunk 2 should start with " fast".
	expectedSecond := " fast engine. We are"
	if chunks[1] != expectedSecond {
		t.Errorf("Expected second chunk %q, got %q", expectedSecond, chunks[1])
	}
}