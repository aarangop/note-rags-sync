// cmd/debug/main.go - create this file for debugging
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	logger.Info("🔧 Debug Mode: Obsidian Sync")
	logger.Infof("📋 Config: %s", cfg.String())

	// Check if vault path exists and list some files
	logger.Infof("📁 Checking vault path: %s", cfg.VaultPath)
	entries, err := os.ReadDir(cfg.VaultPath)
	if err != nil {
		logger.Fatalf("Cannot read vault directory: %v", err)
	}

	logger.Infof("📄 Found %d entries in vault:", len(entries))
	for i, entry := range entries {
		if i >= 5 { // Show first 5 entries
			logger.Infof("   ... and %d more", len(entries)-5)
			break
		}
		logger.Infof("   %s (dir: %t)", entry.Name(), entry.IsDir())
	}

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	w := watcher.New(cfg.VaultPath)

	// Start watcher in a goroutine so we can handle shutdown
	go func() {
		if err := w.Start(); err != nil {
			logger.Fatalf("Failed to start watcher: %v", err)
		}
	}()

	logger.Info("🎯 Watcher started! Try creating/editing .md files in your vault")
	logger.Info("💡 Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-c
	logger.Info("🛑 Shutting down...")
	w.Stop()
	logger.Info("✅ Goodbye!")
}
