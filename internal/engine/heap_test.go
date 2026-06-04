// internal/engine/heap_test.go
package engine

import (
	"testing"
)

func TestResultSetTopK(t *testing.T) {
	// We want to track exactly the top 3 closest items.
	k := 3
	rs := NewResultSet(k)

	// Simulate scanning 5 vectors with various distances.
	// Distances: 5.0, 1.0, 3.0, 2.0, 4.0
	rs.Add(1, 5.0)
	rs.Add(2, 1.0)
	rs.Add(3, 3.0)
	rs.Add(4, 2.0)
	rs.Add(5, 4.0)

	results := rs.Results()

	// 1. Assert the heap strictly maintained the K boundary.
	if len(results) != k {
		t.Fatalf("Expected exactly %d results, got %d", k, len(results))
	}

	// 2. Assert the results are correctly ordered from closest (smallest) to furthest.
	expectedDistances := []float32{1.0, 2.0, 3.0}
	expectedIDs := []uint64{2, 4, 3}

	for i, res := range results {
		if res.Distance != expectedDistances[i] {
			t.Errorf("Index %d: Expected distance %f, got %f", i, expectedDistances[i], res.Distance)
		}
		if res.ID != expectedIDs[i] {
			t.Errorf("Index %d: Expected ID %d, got %d", i, expectedIDs[i], res.ID)
		}
	}
}