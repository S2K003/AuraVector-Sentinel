// internal/engine/bench_test.go
package engine

import (
	"math/rand"
	"testing"
)

// generateSIFTVector creates a mock 128-dimensional vector common in SIFT datasets.
func generateSIFTVector() []float32 {
	v := make([]float32, 128)
	for i := 0; i < 128; i++ {
		v[i] = rand.Float32()
	}
	return v
}

func BenchmarkMathPipelines(b *testing.B) {
	v1 := generateSIFTVector()
	v2 := generateSIFTVector()

	b.Run("1_Baseline_StandardLoop", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CosineSimilarity(v1, v2)
		}
	})

	b.Run("2_Optimized_BCE", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CosineSimilarityBCE(v1, v2)
		}
	})

	b.Run("3_Vectorized_Unrolled", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CosineSimilarityUnrolled(v1, v2)
		}
	})
}