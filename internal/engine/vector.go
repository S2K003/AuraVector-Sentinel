// internal/engine/vector.go
package engine

// Vector represents a single high-dimensional data point.
type Vector struct {
	ID       uint64
	Values   []float32
	Metadata map[string]string // NEW: Key-value tags for Hybrid Search pre-filtering
}