// internal/engine/math_optimized.go
package engine

import "math"

// CosineSimilarityBCE calculates the cosine similarity with Bounds Check Elimination.
func CosineSimilarityBCE(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	// This single line proves to the Go compiler that 'a' is at least as long as 'b'.
	// This completely eliminates runtime bounds checking from the loop below!
	_ = a[len(b)-1]

	var dotProduct, normA, normB float32
	for i := 0; i < len(b); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// SquaredEuclideanBCE calculates distance with Bounds Check Elimination.
func SquaredEuclideanBCE(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	// Explicit bounds check elimination
	_ = a[len(b)-1]

	var distance float32
	for i := 0; i < len(b); i++ {
		diff := a[i] - b[i]
		distance += diff * diff
	}

	return distance
}