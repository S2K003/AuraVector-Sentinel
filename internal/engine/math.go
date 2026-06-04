// internal/engine/math.go
package engine

import "math"

// CosineSimilarity calculates the cosine similarity between two float32 vectors.
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	// We cast to float64 for the Sqrt to utilize standard math libraries, 
	// then back to float32 for our unified schema.
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// SquaredEuclidean calculates the squared Euclidean distance. 
// We skip the final square root calculation for performance, as squared distances 
// sort exactly the same as true Euclidean distances.
func SquaredEuclidean(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	var distance float32
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		distance += diff * diff
	}

	return distance
}