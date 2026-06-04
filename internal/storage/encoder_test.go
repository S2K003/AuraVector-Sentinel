// internal/storage/encoder_test.go
package storage

import (
	"os"
	"testing"

	"auravector/internal/engine"
)

func TestEncodeVector(t *testing.T) {
	v := engine.Vector{
		ID:     42,
		Values: []float32{1.5, 2.5, 3.5},
	}

	encoded := EncodeVector(v, false)
	expectedSize := VectorRecordSize(3) // 1 + 8 + 4 + (4*3) = 25

	if len(encoded) != expectedSize {
		t.Fatalf("Expected byte array of size %d, got %d", expectedSize, len(encoded))
	}

	if encoded[0] != 0 {
		t.Errorf("Expected tombstone flag 0, got %d", encoded[0])
	}
}

func TestWriteAheadLog(t *testing.T) {
	tempFile := "test_data.vdb"
	wal, err := NewWAL(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize WAL: %v", err)
	}
	defer os.Remove(tempFile) // Clean up after test

	v := engine.Vector{ID: 100, Values: []float32{0.1, 0.9}}
	err = wal.Append(v, false)
	if err != nil {
		t.Fatalf("Failed to append vector to WAL: %v", err)
	}
	
	wal.Close()

	info, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("Failed to read WAL file stats: %v", err)
	}

	expectedSize := int64(VectorRecordSize(2))
	if info.Size() != expectedSize {
		t.Errorf("Expected WAL file size %d, got %d", expectedSize, info.Size())
	}
}