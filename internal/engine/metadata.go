// internal/engine/metadata.go
package engine

import "sync"

// MetadataIndex provides exact-match filtering for structured scalar data.
type MetadataIndex struct {
	// store maps: Field -> Value -> Set of allowed Vector IDs
	store map[string]map[string]map[uint64]bool
	mu    sync.RWMutex
}

// NewMetadataIndex initializes a thread-safe scalar index.
func NewMetadataIndex() *MetadataIndex {
	return &MetadataIndex{
		store: make(map[string]map[string]map[uint64]bool),
	}
}

// Add indexes a new set of key-value metadata tags for a specific vector ID.
func (m *MetadataIndex) Add(id uint64, metadata map[string]string) {
	if len(metadata) == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for key, val := range metadata {
		// Initialize maps if they don't exist
		if _, exists := m.store[key]; !exists {
			m.store[key] = make(map[string]map[uint64]bool)
		}
		if _, exists := m.store[key][val]; !exists {
			m.store[key][val] = make(map[uint64]bool)
		}
		
		// Map this exact Key:Value pair to the Vector ID
		m.store[key][val][id] = true
	}
}

// Match returns a set of allowed IDs that match ALL provided filters (Strict AND logic).
// If the filters map is empty, it returns nil, implying "No Restrictions."
func (m *MetadataIndex) Match(filters map[string]string) map[uint64]bool {
	if len(filters) == 0 {
		return nil // No filters applied, allow everything
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var allowed map[uint64]bool
	first := true

	for key, val := range filters {
		matchingIDs := m.store[key][val] // Might be empty/nil

		if first {
			// Initialize the allowed set with the first filter's results
			allowed = make(map[uint64]bool)
			for id := range matchingIDs {
				allowed[id] = true
			}
			first = false
		} else {
			// Perform Set Intersection (AND logic) with subsequent filters
			for id := range allowed {
				if !matchingIDs[id] {
					delete(allowed, id) // If it fails a subsequent filter, remove it
				}
			}
		}

		// Early exit: if intersection ever reaches zero, no need to check further filters
		if len(allowed) == 0 {
			return allowed
		}
	}

	return allowed
}