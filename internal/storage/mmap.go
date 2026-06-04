// internal/storage/mmap.go
package storage

import (
	"os"
	"syscall"
	"unsafe"
)

// MemoryMap represents a read-only memory-mapped file.
type MemoryMap struct {
	data []byte
	addr uintptr
}

// MapRegion loads a file into read-only virtual memory.
// We are using Windows-native syscalls to avoid external C-bindings.
func MapRegion(filepath string) (*MemoryMap, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Safe to close after mapping is established

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() == 0 {
		return &MemoryMap{data: []byte{}}, nil
	}

	// 1. Create a Windows file mapping object
	handle, err := syscall.CreateFileMapping(syscall.Handle(file.Fd()), nil, syscall.PAGE_READONLY, 0, 0, nil)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	// 2. Map the view of the file into our address space
	addr, err := syscall.MapViewOfFile(handle, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 3. Cast the raw memory address into a Go byte slice safely without allocation
	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), int(info.Size()))

	return &MemoryMap{data: data, addr: addr}, nil
}

// Unmap safely releases the virtual memory region back to the OS.
func (m *MemoryMap) Unmap() error {
	if m.addr == 0 {
		return nil
	}
	return syscall.UnmapViewOfFile(m.addr)
}

// ReadVectorZeroAlloc extracts a float32 slice directly from the mapped byte array 
// without triggering the Go garbage collector.
func (m *MemoryMap) ReadVectorZeroAlloc(offset int, dimensions int) []float32 {
	// The raw floats start at offset + 13 bytes (Tombstone 1 + ID 8 + Dims 4)
	floatOffset := offset + 13
	
	// Grab the unsafe pointer to the starting byte of the float array
	ptr := unsafe.Pointer(&m.data[floatOffset])
	
	// Cast the raw memory directly into a Go float32 slice
	// We MUST cast to (*float32) before passing to unsafe.Slice
	return unsafe.Slice((*float32)(ptr), dimensions)
}