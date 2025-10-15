package app

import (
	"fmt"
	"io"
	"os"
)

// IsFileBinary iterates over first 1024 bytes of a file, returning true if the sample contains null bytes or false otherwise
func IsFileBinary(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	sample := make([]byte, 1024)
	bytesRead, err := file.Read(sample)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	for i := 0; i < bytesRead; i++ {
		if sample[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}
