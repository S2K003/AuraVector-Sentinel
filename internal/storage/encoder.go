// internal/storage/encoder.go
package storage

import (
	"encoding/binary"
	"math"
	"os"

	"auravector/internal/engine"
)

// VectorRecordSize calculates the exact byte footprint of a serialized vector.
func VectorRecordSize(dimensions int) int {
	// 1 (Tombstone) + 8 (ID) + 4 (Dimensions) + (4 * Dimensions)
	return 13 + (4 * dimensions)
}

// EncodeVector serializes a Vector into our flat binary byte layout.
func EncodeVector(v engine.Vector, isDeleted bool) []byte {
	dims := len(v.Values)
	size := VectorRecordSize(dims)
	buf := make([]byte, size)

	// [Tombstone Flag (1 Byte)]
	if isDeleted {
		buf[0] = 1
	} else {
		buf[0] = 0
	}

	// [Vector ID (8 Bytes)] - Little Endian
	binary.LittleEndian.PutUint64(buf[1:9], v.ID)

	// [Dimensions (4 Bytes)] - Little Endian
	binary.LittleEndian.PutUint32(buf[9:13], uint32(dims))

	// [Raw Float32 Arrays (4 * N Bytes)]
	offset := 13
	for _, val := range v.Values {
		bits := math.Float32bits(val)
		binary.LittleEndian.PutUint32(buf[offset:offset+4], bits)
		offset += 4
	}

	return buf
}

// WriteAheadLog represents a simple append-only file for durability.
type WriteAheadLog struct {
	file *os.File
}

// NewWAL initializes or opens the append-only log format tracking incoming mutations.
func NewWAL(filepath string) (*WriteAheadLog, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &WriteAheadLog{file: f}, nil
}

// Append writes a serialized vector directly to the log.
func (wal *WriteAheadLog) Append(v engine.Vector, isDeleted bool) error {
	data := EncodeVector(v, isDeleted)
	_, err := wal.file.Write(data)
	return err
}

// Close safely flushes and closes the WAL file.
func (wal *WriteAheadLog) Close() error {
	// Sync ensures data is actually written to the physical disk
	_ = wal.file.Sync()
	return wal.file.Close()
}