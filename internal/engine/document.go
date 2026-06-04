// internal/engine/document.go
package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// DocumentStore maps vector IDs to their original string chunks and persists to disk.
type DocumentStore struct {
	texts    map[uint64]string
	mu       sync.RWMutex
	dataFile *os.File
	MaxID    uint64 // Tracks the highest ID recovered for server coordination
}

type docEntry struct {
	ID   uint64 `json:"id"`
	Text string `json:"text"`
}

// NewDocumentStore initializes a persistent text mapping store.
func NewDocumentStore(filePath string) *DocumentStore {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic("Failed to open document store file: " + err.Error())
	}

	store := &DocumentStore{
		texts:    make(map[uint64]string),
		dataFile: file,
	}

	store.loadFromDisk() // Recover state immediately on boot
	return store
}

// loadFromDisk parses the JSONL file to rebuild the memory map on startup.
func (s *DocumentStore) loadFromDisk() {
	s.dataFile.Seek(0, 0)
	scanner := bufio.NewScanner(s.dataFile)
	
	var recovered int
	for scanner.Scan() {
		var entry docEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err == nil {
			s.texts[entry.ID] = entry.Text
			if entry.ID > s.MaxID {
				s.MaxID = entry.ID
			}
			recovered++
		}
	}
	fmt.Printf("DocumentStore: Recovered %d textual chunks from disk.\n", recovered)
}

// Insert safely adds a new text chunk to memory and disk.
func (s *DocumentStore) Insert(id uint64, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Save to RAM for instant retrieval
	s.texts[id] = text

	// 2. Append to Disk (Write-Ahead Log style)
	entry, _ := json.Marshal(docEntry{ID: id, Text: text})
	s.dataFile.Write(append(entry, '\n'))
}

// Get safely retrieves a text chunk by its vector ID.
func (s *DocumentStore) Get(id uint64) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	text, exists := s.texts[id]
	return text, exists
}