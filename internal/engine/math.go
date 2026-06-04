// internal/engine/math.go
package engine

import (
	"fmt"
	"math"
)

// EuclideanDistance calculates the straight-line distance between two 768-D vectors.
// The larger the result, the more semantically different the two web requests are.
func EuclideanDistance(v1, v2 []float32) (float64, error) {
	if len(v1) != len(v2) {
		// FIX: Return 0 for the float, and a proper Go error for the error interface
		return 0, fmt.Errorf("dimension mismatch: vector 1 is %d-D, vector 2 is %d-D", len(v1), len(v2))
	}

	var sum float64
	for i := 0; i < len(v1); i++ {
		diff := float64(v1[i] - v2[i])
		sum += diff * diff
	}

	return math.Sqrt(sum), nil
}