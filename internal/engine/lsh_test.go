// internal/engine/lsh_test.go
package engine

import (
	"testing"
)

func TestLSHHash(t *testing.T) {
	// Initialize LSH with 3 dimensions, 4 hyperplanes, and a fixed seed of 42 for determinism
	lsh := NewLSH(3, 4, 42)

	// Two vectors that are extremely close to one another in space
	v1 := []float32{1.0, 1.0, 1.0}
	v2 := []float32{1.0, 1.1, 0.9} 

	// A vector pointing in the complete opposite direction
	v3 := []float32{-1.0, -1.0, -1.0} 

	sig1 := lsh.Hash(v1)
	sig2 := lsh.Hash(v2)
	sig3 := lsh.Hash(v3)

	// Expectation 1: Close vectors should fall into the exact same spatial bucket (same signature)
	if sig1 != sig2 {
		t.Errorf("Expected identical signatures for close vectors, got %d and %d", sig1, sig2)
	}

	// Expectation 2: Opposing vectors should land on the opposite side of the hyperplanes
	if sig1 == sig3 {
		t.Errorf("Expected different signatures for opposing vectors, got %d", sig3)
	}
}