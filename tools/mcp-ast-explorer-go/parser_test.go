package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndExtract(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy python file
	pyFile := filepath.Join(tempDir, "test.py")
	pyCode := `
def hello_world(name: str) -> str:
    """This is a docstring."""
    return f"Hello, {name}"

class MyClass:
    def method(self):
        pass
`

	err := os.WriteFile(pyFile, []byte(pyCode), 0644)
	if err != nil {
		t.Fatalf("failed to write python file: %v", err)
	}

	relPath, nodes := parseAndExtract(pyFile, tempDir, "python")
	if relPath == "" || len(nodes) == 0 {
		t.Fatalf("expected nodes, got none")
	}

	foundHello := false
	foundClass := false
	for _, n := range nodes {
		if n.Name == "hello_world" {
			foundHello = true
			if n.Doc != "This is a docstring." {
				t.Errorf("Expected docstring 'This is a docstring.', got '%s'", n.Doc)
			}
			if n.Signature != "(name: str) -> str" {
				t.Errorf("Expected signature '(name: str) -> str', got '%s'", n.Signature)
			}
		}
		if n.Name == "MyClass" {
			foundClass = true
		}
	}

	if !foundHello || !foundClass {
		t.Errorf("Missing expected nodes in Python parser")
	}

	// Create a dummy go file
	goFile := filepath.Join(tempDir, "test.go")
	goCode := `package main

// MyGoFunc does something
func MyGoFunc(a int) bool {
	return true
}
`
	os.WriteFile(goFile, []byte(goCode), 0644)

	_, goNodes := parseAndExtract(goFile, tempDir, "go")
	if len(goNodes) == 0 {
		t.Fatalf("expected go nodes, got none")
	}

	foundGo := false
	for _, n := range goNodes {
		if n.Name == "MyGoFunc" {
			foundGo = true
			if n.Doc != "MyGoFunc does something" {
				t.Errorf("Expected go docstring 'MyGoFunc does something', got '%s'", n.Doc)
			}
			if n.Signature != "(a int) bool" {
				t.Errorf("Expected signature '(a int) bool', got '%s'", n.Signature)
			}
		}
	}

	if !foundGo {
		t.Errorf("Missing expected nodes in Go parser")
	}
}
