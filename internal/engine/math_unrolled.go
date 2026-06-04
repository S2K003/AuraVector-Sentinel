// internal/engine/math_unrolled.go
package engine

import "math"

// CosineSimilarityUnrolled calculates cosine similarity using 4-wide loop unrolling
// to encourage compiler SIMD (Single Instruction, Multiple Data) vectorization.
func CosineSimilarityUnrolled(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	// Prove bounds to the compiler (BCE)
	_ = a[len(b)-1]

	var dotProduct, normA, normB float32
	
	// Process the arrays in 4-wide chunks
	i := 0
	for ; i <= len(b)-4; i += 4 {
		// Unrolled dot product
		dotProduct += a[i]*b[i] + a[i+1]*b[i+1] + a[i+2]*b[i+2] + a[i+3]*b[i+3]
		
		// Unrolled Euclidean norms
		normA += a[i]*a[i] + a[i+1]*a[i+1] + a[i+2]*a[i+2] + a[i+3]*a[i+3]
		normB += b[i]*b[i] + b[i+1]*b[i+1] + b[i+2]*b[i+2] + b[i+3]*b[i+3]
	}

	// Clean up any remaining elements if the length isn't perfectly divisible by 4
	for ; i < len(b); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}