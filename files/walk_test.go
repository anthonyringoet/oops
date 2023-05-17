// files/files_test.go
package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestWalkFiles(t *testing.T) {
	// Create a temporary directory for testing
	root, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(root)

	// Create some files in the directory
	for i := 0; i < 10; i++ {
		file, err := ioutil.TempFile(root, "test")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		file.Close()
	}

	walker := NewFileWalker(root)
	files, err := walker.WalkFiles()
	if err != nil {
		t.Fatalf("Failed to walk files: %v", err)
	}

	if len(files) != 10 {
		t.Errorf("Expected 10 files, got %d", len(files))
	}
}
