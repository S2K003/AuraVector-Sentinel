// internal/engine/math_test.go
package engine

import (
	"math"
	"testing"
)

// almostEqual is a helper to handle floating-point precision comparisons
func almostEqual(a, b, tolerance float32) bool {
	return float32(math.Abs(float64(a-b))) <= tolerance
}

func TestCosineSimilarity(t *testing.T) {
	v1 := []float32{1.0, 0.0, 0.0}
	v2 := []float32{0.0, 1.0, 0.0}
	v3 := []float32{1.0, 0.0, 0.0}

	// Orthogonal vectors should have 0 similarity
	if sim := CosineSimilarity(v1, v2); !almostEqual(sim, 0.0, 1e-5) {
		t.Errorf("Expected 0.0, got %f", sim)
	}

	// Identical vectors should have 1 similarity
	if sim := CosineSimilarity(v1, v3); !almostEqual(sim, 1.0, 1e-5) {
		t.Errorf("Expected 1.0, got %f", sim)
	}
}

func TestSquaredEuclidean(t *testing.T) {
	v1 := []float32{0.0, 0.0}
	v2 := []float32{3.0, 4.0}

	// Squared distance should be 3^2 + 4^2 = 9 + 16 = 25
	expected := float32(25.0)
	if dist := SquaredEuclidean(v1, v2); !almostEqual(dist, expected, 1e-5) {
		t.Errorf("Expected %f, got %f", expected, dist)
	}
}