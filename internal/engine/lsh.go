// internal/engine/lsh.go
package engine

import (
	"math/rand"
	"time"
)

// LSH represents a Locality-Sensitive Hashing configuration.
type LSH struct {
	Hyperplanes [][]float32
	NumPlanes   int
}

// NewLSH initializes K random hyperplanes for chunking spatial regions.
// We restrict numPlanes to a maximum of 32 to fit within our uint32 bitmask.
func NewLSH(dimensions int, numPlanes int, seed int64) *LSH {
	if numPlanes > 32 {
		numPlanes = 32
	}

	// Use a fixed seed for deterministic testing, or current time for production
	r := rand.New(rand.NewSource(seed))
	if seed == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	planes := make([][]float32, numPlanes)
	for i := 0; i < numPlanes; i++ {
		plane := make([]float32, dimensions)
		for j := 0; j < dimensions; j++ {
			// Generate random values between -1.0 and 1.0 to create the normal vector
			plane[j] = r.Float32()*2.0 - 1.0
		}
		planes[i] = plane
	}

	return &LSH{
		Hyperplanes: planes,
		NumPlanes:   numPlanes,
	}
}

// Hash evaluates the dot product of a vector against all hyperplanes.
// It returns a uint32 bitmask signature representing the LSH region.
func (lsh *LSH) Hash(vector []float32) uint32 {
	var signature uint32 = 0

	for i := 0; i < lsh.NumPlanes; i++ {
		var dotProduct float32 = 0
		for j := 0; j < len(vector); j++ {
			dotProduct += vector[j] * lsh.Hyperplanes[i][j]
		}

		// If the spatial coordinate calculates greater than 0, flip the i-th bit to 1
		if dotProduct > 0 {
			signature |= (1 << i)
		}
	}

	return signature
}