package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	path := "/tmp/test"
	w := New(path)

	if w.path != path {
		t.Errorf("Expected path %s, got %s", path, w.path)
	}
}

func TestWatcherWithTempDir(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up

	w := New(tmpDir)

	// Start watcher in goroutine since it blocks
	go func() {
		if err := w.Start(); err != nil {
			t.Errorf("Failed to start watcher: %v", err)
		}
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// In a real test, you'd verify the watcher detected the change
	// For now, this just ensures no crashes
}
