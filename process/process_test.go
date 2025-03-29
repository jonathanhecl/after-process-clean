//go:build windows

package process

import (
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	// Test existing file
	if !Exists(tempFilePath) {
		t.Errorf("Exists() returned false for existing file %s", tempFilePath)
	}

	// Test non-existing file
	nonExistingPath := tempFilePath + "_nonexistent"
	if Exists(nonExistingPath) {
		t.Errorf("Exists() returned true for non-existing file %s", nonExistingPath)
	}
}

func TestList(t *testing.T) {
	// This is a basic smoke test to ensure the List function doesn't crash
	// and returns some processes on Windows
	processes := List()
	
	// On a running Windows system, we should have at least a few processes
	if len(processes) == 0 {
		t.Error("List() returned empty process list, expected at least some system processes")
	}

	// Check that at least some processes have non-empty paths
	hasNonEmptyPath := false
	for _, p := range processes {
		if p.Path != "" {
			hasNonEmptyPath = true
			break
		}
	}

	if !hasNonEmptyPath {
		t.Error("Expected at least one process to have a non-empty path")
	}
}
