// internal/engine/vector_test.go
package engine

import (
	"testing"
)

func TestVectorInitialization(t *testing.T) {
	v := Vector{
		ID:     1,
		Values: []float32{0.1, 0.2, 0.3},
	}

	if v.ID != 1 {
		t.Errorf("Expected ID 1, got %d", v.ID)
	}

	if len(v.Values) != 3 {
		t.Errorf("Expected Values length 3, got %d", len(v.Values))
	}
}