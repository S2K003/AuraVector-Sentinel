// internal/engine/index.go
package engine

import "sync"

// InvertedIndex manages the mapping of LSH signatures to a list of Vector IDs.
type InvertedIndex struct {
	mu      sync.RWMutex
	buckets map[uint32][]uint64
}

// NewInvertedIndex initializes an empty protected index.
func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		buckets: make(map[uint32][]uint64),
	}
}

// Add maps a vector ID to its corresponding LSH signature bucket.
func (idx *InvertedIndex) Add(signature uint32, id uint64) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.buckets[signature] = append(idx.buckets[signature], id)
}

// GetExactCandidates retrieves vector IDs from the identical spatial region.
func (idx *InvertedIndex) GetExactCandidates(signature uint32) []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	
	candidates := make([]uint64, len(idx.buckets[signature]))
	copy(candidates, idx.buckets[signature])
	return candidates
}

// GetNeighboringCandidates retrieves candidates from buckets that differ 
// by exactly 1 bit (Hamming distance 1) to expand the ANN search radius.
func (idx *InvertedIndex) GetNeighboringCandidates(signature uint32, numPlanes int) []uint64 {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var candidates []uint64
	// Flip each bit one by one to find neighboring spatial buckets
	for i := 0; i < numPlanes; i++ {
		neighborSignature := signature ^ (1 << i)
		candidates = append(candidates, idx.buckets[neighborSignature]...)
	}
	return candidates
}