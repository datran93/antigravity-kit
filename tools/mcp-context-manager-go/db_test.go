package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDBConnection(t *testing.T) {
	tempDir := t.TempDir()
	db, err := GetDBConnection(tempDir)
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Fatal("Expected db connection, got nil")
	}

	// Verify schema was created
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='checkpoints'").Scan(&name)
	if err != nil {
		t.Fatalf("Table checkpoints not found: %v", err)
	}
}

func TestCheckpointOperations(t *testing.T) {
	tempDir := t.TempDir()
	taskID := "test-task-1"

	// Create
	res, err := InitializeTaskPlan(tempDir, taskID, "Test description", []string{"Step 1", "Step 2"})
	if err != nil {
		t.Fatalf("Failed to initialize task: %v", err)
	}
	if res == "" {
		t.Error("Expected result string")
	}

	// Complete step
	res, err = CompleteTaskStep(tempDir, taskID, "Step 1", []string{"file1.go"}, "Did step 1")
	if err != nil {
		t.Fatalf("Failed to complete step: %v", err)
	}

	// Add step
	_, err = AddTaskStep(tempDir, taskID, "Step 3")
	if err != nil {
		t.Fatalf("Failed to add step: %v", err)
	}

	// Load checkpoint
	out, err := LoadCheckpoint(tempDir, taskID)
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	// Verify progress.md exists
	mdPath := filepath.Join(tempDir, "progress.md")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		t.Error("progress.md was not created")
	}

	if out == "" {
		t.Error("Expected checkpoint string")
	}
}

func TestGovernanceOperations(t *testing.T) {
	tempDir := t.TempDir()
	taskID := "test-task-2"

	// Intent declaration
	_, err := DeclareIntent(tempDir, taskID, "Test intent", []string{"fileA.go", "fileB.go"})
	if err != nil {
		t.Fatalf("Failed to declare intent: %v", err)
	}

	// Check valid intent
	res, err := CheckIntentLock(tempDir, taskID, "fileA.go")
	if err != nil {
		t.Fatalf("Failed to check intent: %v", err)
	}
	if res == "" {
		t.Error("Expected check intent result")
	}

	// Annotate
	_, err = AnnotateFile(tempDir, "fileA.go", "Important ghost context")
	if err != nil {
		t.Fatalf("Failed to annotate file: %v", err)
	}

	// Record drift
	res, err = RecordFailure(tempDir, taskID)
	if err != nil || res == "" {
		t.Fatalf("Failed to record failure")
	}

	// Clear drift
	_, err = ClearDrift(tempDir, taskID)
	if err != nil {
		t.Fatalf("Failed to clear drift")
	}

	// Manage anchors
	_, err = ManageAnchors(tempDir, "set", "global_rule", "val", "rule1")
	if err != nil {
		t.Fatalf("Failed to manage anchors")
	}
}
