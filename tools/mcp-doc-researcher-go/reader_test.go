package main

import (
	"strings"
	"testing"
)

func TestReadLocalDoc(t *testing.T) {
	content, err := ReadLocalDoc("test_large.docx")
	if err != nil {
		t.Fatalf("Failed to read test_large.docx: %v", err)
	}
	if !strings.Contains(content, "x") {
		t.Errorf("Expected content to contain 'x', got parsed text length: %d", len(content))
	} else {
		t.Logf("Successfully extracted %d characters", len(content))
	}
}
