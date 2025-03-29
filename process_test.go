package main

import (
	"testing"

	"github.com/afterprocessclean/process"
)

// Mock process list for testing
func mockProcessList() []process.ProcessStruct {
	return []process.ProcessStruct{
		{PID: 1, Filename: "process1.exe", Path: "C:\\path\\to\\process1.exe"},
		{PID: 2, Filename: "process2.exe", Path: "C:\\path\\to\\process2.exe"},
	}
}

func TestControlStructAddProcess(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Test adding a process
	ctrl.addProcess("C:\\path\\to\\test.exe", "12345", true)

	if len(ctrl.processes) != 1 {
		t.Errorf("Expected 1 process after adding, got %d", len(ctrl.processes))
	}

	p := ctrl.processes[0]
	if p.Path != "C:\\path\\to\\test.exe" {
		t.Errorf("Expected path 'C:\\path\\to\\test.exe', got '%s'", p.Path)
	}
	if p.CRC32 != "12345" {
		t.Errorf("Expected CRC32 '12345', got '%s'", p.CRC32)
	}
	if !p.Before {
		t.Errorf("Expected Before to be true, got false")
	}
}

func TestControlStructGetProcess(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Add a test process
	testPath := "C:\\path\\to\\test.exe"
	ctrl.addProcess(testPath, "12345", true)

	// Test getting an existing process
	p := ctrl.getProcess(testPath)
	if p == nil {
		t.Errorf("Expected to find process with path '%s', got nil", testPath)
	}

	// Test getting a non-existent process
	p = ctrl.getProcess("nonexistent.exe")
	if p != nil {
		t.Errorf("Expected nil for non-existent process, got %v", p)
	}
}

func TestControlStructRemoveProcess(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Add test processes
	testPath1 := "C:\\path\\to\\test1.exe"
	testPath2 := "C:\\path\\to\\test2.exe"
	ctrl.addProcess(testPath1, "12345", true)
	ctrl.addProcess(testPath2, "67890", true)

	// Verify initial state
	if len(ctrl.processes) != 2 {
		t.Fatalf("Expected 2 processes before removal, got %d", len(ctrl.processes))
	}

	// Test removing a process
	ctrl.removeProcess(testPath1)

	// Verify process was removed
	if len(ctrl.processes) != 1 {
		t.Errorf("Expected 1 process after removal, got %d", len(ctrl.processes))
	}
	if ctrl.processes[0].Path != testPath2 {
		t.Errorf("Expected remaining process to be '%s', got '%s'", testPath2, ctrl.processes[0].Path)
	}
}

func TestControlStructUpdateProcess(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Add a test process
	testPath := "C:\\path\\to\\test.exe"
	initialCRC := "12345"
	ctrl.addProcess(testPath, initialCRC, true)

	// Update the process CRC32
	newCRC := "67890"
	ctrl.updateProcess(testPath, newCRC)

	// Verify the CRC32 was updated
	p := ctrl.getProcess(testPath)
	if p == nil {
		t.Fatalf("Process not found after update")
	}
	if p.CRC32 != newCRC {
		t.Errorf("Expected updated CRC32 '%s', got '%s'", newCRC, p.CRC32)
	}
}

func TestControlStructAfterList(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Add test processes with different Before values
	ctrl.addProcess("C:\\path\\to\\before.exe", "12345", true)
	ctrl.addProcess("C:\\path\\to\\after.exe", "67890", false)

	// Get the list of processes added after initialization
	afterList := ctrl.AfterList()

	// Verify only processes with Before=false are returned
	if len(afterList) != 1 {
		t.Errorf("Expected 1 process in AfterList, got %d", len(afterList))
	}
	if len(afterList) > 0 && afterList[0].Path != "C:\\path\\to\\after.exe" {
		t.Errorf("Expected path 'C:\\path\\to\\after.exe', got '%s'", afterList[0].Path)
	}
}

func TestControlStructUpdateList(t *testing.T) {
	// Initialize a new controlStruct for testing
	ctrl := controlStruct{
		processes: []*processStruct{},
	}

	// Create initial process list
	initialList := []process.ProcessStruct{
		{PID: 1, Filename: "process1.exe", Path: "C:\\path\\to\\process1.exe"},
		{PID: 2, Filename: "process2.exe", Path: "C:\\path\\to\\process2.exe"},
	}

	// Update with initial list
	ctrl.UpdateList(initialList, true)

	// Verify processes were added
	if len(ctrl.processes) != 2 {
		t.Fatalf("Expected 2 processes after initial update, got %d", len(ctrl.processes))
	}

	// Create updated process list (one removed, one added)
	updatedList := []process.ProcessStruct{
		{PID: 2, Filename: "process2.exe", Path: "C:\\path\\to\\process2.exe"},
		{PID: 3, Filename: "process3.exe", Path: "C:\\path\\to\\process3.exe"},
	}

	// Update with new list
	ctrl.UpdateList(updatedList, false)

	// Verify processes were updated correctly
	if len(ctrl.processes) != 2 {
		t.Fatalf("Expected 2 processes after update, got %d", len(ctrl.processes))
	}

	// Check that process1 was removed and process3 was added
	paths := make(map[string]bool)
	for _, p := range ctrl.processes {
		paths[p.Path] = true
	}

	if paths["C:\\path\\to\\process1.exe"] {
		t.Error("Expected process1.exe to be removed, but it's still present")
	}
	if !paths["C:\\path\\to\\process3.exe"] {
		t.Error("Expected process3.exe to be added, but it's not present")
	}
}
