// internal/engine/metadata_test.go
package engine

import "testing"

func TestMetadataIndex(t *testing.T) {
	idx := NewMetadataIndex()

	// Seed dummy data
	idx.Add(1, map[string]string{"author": "Thorne", "category": "science"})
	idx.Add(2, map[string]string{"author": "Thorne", "category": "fiction"})
	idx.Add(3, map[string]string{"author": "Squawks", "category": "science"})

	// Scenario 1: Single Filter match
	res1 := idx.Match(map[string]string{"author": "Thorne"})
	if len(res1) != 2 || !res1[1] || !res1[2] {
		t.Errorf("Expected IDs 1 and 2, got %v", res1)
	}

	// Scenario 2: Multiple Filters (AND logic)
	res2 := idx.Match(map[string]string{"author": "Thorne", "category": "science"})
	if len(res2) != 1 || !res2[1] {
		t.Errorf("Expected strictly ID 1, got %v", res2)
	}

	// Scenario 3: Impossible overlap (Empty result)
	res3 := idx.Match(map[string]string{"author": "Squawks", "category": "fiction"})
	if len(res3) != 0 {
		t.Errorf("Expected empty result, got %v", res3)
	}
	
	// Scenario 4: No filters provided
	res4 := idx.Match(map[string]string{})
	if res4 != nil {
		t.Errorf("Expected nil allowed map, got %v", res4)
	}
}