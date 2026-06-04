// internal/engine/heap.go
package engine

import "container/heap"

// VectorDistance represents a scored vector during a search.
type VectorDistance struct {
	ID       uint64
	Distance float32
}

// maxHeap implements container/heap for VectorDistance.
// We use a max-heap so the vector with the LARGEST distance (worst match) is at the root.
type maxHeap []VectorDistance

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].Distance > h[j].Distance } // Max-heap: >
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x any) {
	*h = append(*h, x.(VectorDistance))
}

func (h *maxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ResultSet wraps the maxHeap to strictly enforce the top-K boundary.
type ResultSet struct {
	h maxHeap
	K int
}

// NewResultSet initializes a bounded top-K result tracker.
func NewResultSet(k int) *ResultSet {
	rs := &ResultSet{
		h: make(maxHeap, 0, k),
		K: k,
	}
	heap.Init(&rs.h)
	return rs
}

// Add inserts a new vector distance into the set.
// If the set is full, it evicts the worst match if the new distance is strictly better.
func (rs *ResultSet) Add(id uint64, distance float32) {
	if rs.h.Len() < rs.K {
		heap.Push(&rs.h, VectorDistance{ID: id, Distance: distance})
	} else if distance < rs.h[0].Distance {
		// Since it's a max-heap, h[0] is always the maximum distance (worst match).
		// If our new distance is smaller (better), we pop the worst and push the new.
		heap.Pop(&rs.h)
		heap.Push(&rs.h, VectorDistance{ID: id, Distance: distance})
	}
}

// Results returns the sorted top-K matches (closest first).
func (rs *ResultSet) Results() []VectorDistance {
	res := make([]VectorDistance, rs.h.Len())
	// Pop all elements to sort them. 
	// Because it's a max-heap, popping yields the largest (worst) distances first.
	// We populate the result array backwards so the closest match is at index 0.
	for i := len(res) - 1; i >= 0; i-- {
		res[i] = heap.Pop(&rs.h).(VectorDistance)
	}
	return res
}