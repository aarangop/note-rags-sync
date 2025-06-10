package loaders

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aarangop/obsidian-sync/internal/models"
)

func LoadFile(path string) (models.File, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".md":
		return LoadMarkdownFile(path)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", path)
	}
}

func LoadMarkdownFile(path string) (models.File, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}

	// Read file content
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read files %s: %w", path, err)
	}

	mdFile := &models.MarkdownFile{
		Path:     path,
		Content:  content,
		Checksum: calculateSHA256(content),
	}

	return mdFile, nil
}

func calculateSHA256(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}
