// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "obsidian_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Clean up temp directory after test
	defer os.RemoveAll(tempDir)

	// Create a test markdown file to make it look like a real vault
	testFile := filepath.Join(tempDir, "test.md")
	err = os.WriteFile(testFile, []byte("# Test Note"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set test environment variables
	os.Setenv("VAULT_PATH", tempDir)
	os.Setenv("S3_BUCKET", "test-bucket")
	os.Setenv("APP_VERSION", "test-version")

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("VAULT_PATH")
		os.Unsetenv("S3_BUCKET")
		os.Unsetenv("APP_VERSION")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Version != "test-version" {
		t.Errorf("Expected version 'test-version', got '%s'", cfg.Version)
	}

	if cfg.S3Bucket != "test-bucket" {
		t.Errorf("Expected S3 bucket 'test-bucket', got '%s'", cfg.S3Bucket)
	}

	// Test that the temp directory exists and is accessible
	if cfg.VaultPath != tempDir {
		t.Errorf("Expected vault path '%s', got '%s'", tempDir, cfg.VaultPath)
	}
}
