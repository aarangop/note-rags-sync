package main

import (
	"fmt"
	"os"

	"github.com/aarangop/obsidian-sync/internal/config"
	"github.com/aarangop/obsidian-sync/internal/logger"
	"github.com/aarangop/obsidian-sync/internal/watcher"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logConfig := logger.Config{
		LogLevel:      cfg.LogLevel,
		LogFile:       cfg.LogFile,
		MaxSize:       100, // 100 MB per file
		MaxBackups:    10,  // Keep 10 old files
		MaxAge:        30,  // 30 days
		Compress:      true,
		ConsoleOutput: true,
	}
	logger.Initialize(logConfig)

	logger.Infof("Obsidian Sync v%s", cfg.Version)
	logger.Infof("Configuration loaded %s", cfg.String())

	// Create and start watcher
	w := watcher.New(cfg.VaultPath)

	if err := w.Start(); err != nil {
		logger.Fatalf("Failed to start watcher: %v", err)
	}
}
