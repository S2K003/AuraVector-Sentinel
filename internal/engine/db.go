// internal/engine/db.go
package engine

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

// DB wraps the execution components and coordinates thread-safe operations.
type DB struct {
	mu       sync.RWMutex
	Index    *InvertedIndex
	LSH      *LSH
	Metadata *MetadataIndex // NEW: Structured scalar index
	walFile  *os.File
}

// NewDB initializes the thread-safe database coordinator.
func NewDB(dimensions int, numPlanes int, seed int64) *DB {
	return &DB{
		Index:    NewInvertedIndex(),
		LSH:      NewLSH(dimensions, numPlanes, seed),
		Metadata: NewMetadataIndex(), // NEW
	}
}

// Insert safely adds a vector mapping to the index and writes it to disk.
func (db *DB) Insert(v Vector) {
	db.mu.Lock()
	defer db.mu.Unlock()

	signature := db.LSH.Hash(v.Values)
	db.Index.Add(signature, v.ID)

	// NEW: Index the structured scalar data
	db.Metadata.Add(v.ID, v.Metadata)

	// Write-Ahead Log to persist the vector to disk permanently
	if db.walFile != nil {
		tombstone := []byte{0}
		idBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(idBytes, v.ID)

		dims := uint32(len(v.Values))
		dimBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(dimBytes, dims)

		floatBytes := make([]byte, 4*dims)
		for i, val := range v.Values {
			bits := math.Float32bits(val)
			binary.LittleEndian.PutUint32(floatBytes[i*4 : (i+1)*4], bits)
		}

		db.walFile.Write(tombstone)
		db.walFile.Write(idBytes)
		db.walFile.Write(dimBytes)
		db.walFile.Write(floatBytes)
	}
}

// Search isolates identical and neighboring bucket candidates with pre-filtering.
func (db *DB) Search(query []float32, filters map[string]string) []uint64 {
	// 1. Pre-filter allowed IDs based on metadata
	allowedIDs := db.Metadata.Match(filters)

	// 2. Early Exit: If filters were applied but NO documents matched, skip the math entirely!
	if filters != nil && len(filters) > 0 && len(allowedIDs) == 0 {
		return []uint64{}
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	signature := db.LSH.Hash(query)

	exact := db.Index.GetExactCandidates(signature)
	neighbors := db.Index.GetNeighboringCandidates(signature, db.LSH.NumPlanes)

	candidates := append(exact, neighbors...)

	// 3. Apply the pre-filter mask
	if allowedIDs == nil {
		return candidates // No filters applied, allow all LSH candidates
	}

	var filteredCandidates []uint64
	for _, id := range candidates {
		if allowedIDs[id] { // Only keep candidates that passed the metadata check
			filteredCandidates = append(filteredCandidates, id)
		}
	}

	return filteredCandidates
}

// Recover scans a flat binary .vdb file to rebuild the LSH spatial index in memory.
func (db *DB) Recover(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	db.walFile = file
	db.walFile.Seek(0, 0)

	var recoveredCount int

	for {
		tombstone := make([]byte, 1)
		if _, err := io.ReadFull(db.walFile, tombstone); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		idBytes := make([]byte, 8)
		if _, err := io.ReadFull(db.walFile, idBytes); err != nil {
			return err
		}
		id := binary.LittleEndian.Uint64(idBytes)

		dimBytes := make([]byte, 4)
		if _, err := io.ReadFull(db.walFile, dimBytes); err != nil {
			return err
		}
		dims := binary.LittleEndian.Uint32(dimBytes)

		floatBytes := make([]byte, 4*dims)
		if _, err := io.ReadFull(db.walFile, floatBytes); err != nil {
			return err
		}

		if tombstone[0] == 1 {
			continue
		}

		values := make([]float32, dims)
		for i := uint32(0); i < dims; i++ {
			bits := binary.LittleEndian.Uint32(floatBytes[i*4 : (i+1)*4])
			values[i] = math.Float32frombits(bits)
		}

		signature := db.LSH.Hash(values)
		db.Index.Add(signature, id)
		recoveredCount++
	}

	fmt.Printf("VectorDB: Successfully recovered and re-indexed %d vectors from disk.\n", recoveredCount)
	return nil
}