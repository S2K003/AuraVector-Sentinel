// internal/storage/mmap_test.go
package storage

import (
	"os"
	"testing"
	"auravector/internal/engine"
)

// setupBenchmarkDB is a helper to write a temporary file for reading.
func setupBenchmarkDB(b *testing.B) string {
	tempFile := "bench_data.vdb"
	wal, err := NewWAL(tempFile)
	if err != nil {
		b.Fatal(err)
	}
	
	v := engine.Vector{ID: 99, Values: []float32{1.1, 2.2, 3.3, 4.4}}
	_ = wal.Append(v, false)
	wal.Close()
	
	return tempFile
}

func BenchmarkZeroAllocationRead(b *testing.B) {
	dbPath := setupBenchmarkDB(b)
	defer os.Remove(dbPath)

	mmapRegion, err := MapRegion(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer mmapRegion.Unmap()

	// Reset timer to strictly measure the zero-allocation read operation
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Read a 4-dimensional vector from offset 0
		floatSlice := mmapRegion.ReadVectorZeroAlloc(0, 4)
		
		// Prevent the Go compiler from optimizing the read away
		_ = floatSlice[0]
	}
}