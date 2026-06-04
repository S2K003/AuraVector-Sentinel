// internal/pipeline/chunker.go
package pipeline

// ChunkText splits a string into chunks of a specific rune length,
// with a specified overlap to preserve semantic context across boundaries.
func ChunkText(text string, chunkSize int, overlap int) []string {
	// Cast to runes to safely handle multi-byte characters (like emojis or non-ASCII text)
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	
	// Safety constraints
	if chunkSize <= 0 {
		chunkSize = 250
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 2 
	}

	var chunks []string
	step := chunkSize - overlap

	for i := 0; i < len(runes); i += step {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		
		// If we've reached the end of the text, break the loop
		if end == len(runes) {
			break
		}
	}
	return chunks
}