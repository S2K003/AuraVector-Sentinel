// internal/engine/db_test.go
package engine

import (
	"sync"
	"testing"
)

func TestDBConcurrentAccess(t *testing.T) {
	// Initialize a high-dimensional configuration
	db := NewDB(128, 8, 42)
	var wg sync.WaitGroup

	// Simulate 100 concurrent workers: 50 writing data, 50 reading data simultaneously
	workerCount := 50

	// Launch Writers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			
			// Generate a dummy vector
			v := Vector{ID: id, Values: make([]float32, 128)}
			v.Values[0] = float32(id) // Give it some variance
			
			db.Insert(v)
		}(uint64(i))
	}

	// Launch Readers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			query := make([]float32, 128)
			query[0] = float32(id)
			
			_ = db.Search(query)
		}(i)
	}

	// Wait for all 100 goroutines to finish execution
	wg.Wait()
}