// internal/engine/index_test.go
package engine

import (
	"testing"
)

func TestInvertedIndex(t *testing.T) {
	idx := NewInvertedIndex()
	
	// Simulating a signature of 5 (Binary: 0101)
	var sigExact uint32 = 5 
	
	idx.Add(sigExact, 100)
	idx.Add(sigExact, 101)
	
	// Add to a neighboring bucket by flipping bit 0 (0101 -> 0100 = 4)
	var sigNeighbor uint32 = 4
	idx.Add(sigNeighbor, 200)

	// Test 1: Exact match isolation
	exact := idx.GetExactCandidates(sigExact)
	if len(exact) != 2 {
		t.Errorf("Expected 2 exact candidates, got %d", len(exact))
	}

	// Test 2: Neighboring candidate retrieval
	// Assuming 3 hyperplane dimensions, checking 1-bit flips
	neighbors := idx.GetNeighboringCandidates(sigExact, 3)
	
	// sigExact (5) flips with 3 planes:
	// bit 0: 5 ^ 1 = 4 (sigNeighbor, contains ID 200)
	// bit 1: 5 ^ 2 = 7 (empty)
	// bit 2: 5 ^ 4 = 1 (empty)
	
	if len(neighbors) != 1 || neighbors[0] != 200 {
		t.Errorf("Expected 1 neighboring candidate with ID 200, got %v", neighbors)
	}
}