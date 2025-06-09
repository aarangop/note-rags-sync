package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aarangop/obsidian-sync/internal/logger"
	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a directory for file system events.
// It wraps the fsnotify.Watcher to provide a higher-level interface
// for watching file system changes in a specified path.
type Watcher struct {
	path      string
	fsWatcher *fsnotify.Watcher
	done      chan bool

	eventBuffer   map[string]*fileEvent
	debounceTimer *time.Timer
}

type fileEvent struct {
	path       string
	isNew      bool
	isModified bool
	isDeleted  bool
	lastSeen   time.Time
}

func New(path string) *Watcher {
	return &Watcher{
		path:        path,
		done:        make(chan bool), // Create a channel for clean shutdown
		eventBuffer: make(map[string]*fileEvent),
	}
}

// Start initiates the file watching process.
// It creates a new fsnotify watcher, adds the target directory and all its subdirectories recursively,
// and launches a goroutine to handle file system events.
// The method blocks until the watcher's done channel receives a signal.
//
// Returns an error if creating the watcher or adding directories fails.
func (w *Watcher) Start() error {
	var err error
	w.fsWatcher, err = fsnotify.NewWatcher()

	if err != nil {
		logger.Errorf("⚠️ Failed to create file watcher: %v", err)
		return fmt.Errorf("failed to create file watcher: %v", err)
	}

	// Add existing directories recursively
	err = w.addRecursive(w.path)

	if err != nil {
		logger.Errorf("⚠️ Failed to add directories: %v", err)
		return fmt.Errorf("failed to add directories: %v", err)
	}
	// `go` keyword starts a 'goroutine', a lightweight thread
	go w.watch()

	logger.Infof("🔍Watching for %s for changes...", w.path)

	// Wait for done signal instead of blocking forever
	<-w.done

	return nil
}

// watch starts an infinite monitoring loop for the directory being watched.
// It processes two types of channel events:
//  1. File events: Filters for markdown (.md) files and prints the operation and file name.
//     TODO: Will eventually call an HTTP endpoint to process these events.
//  2. Error events: Logs any errors that occur during watching but continues monitoring.
//
// The function exits when either channel is closed (which happens when the watcher is closed).
func (w *Watcher) watch() {
	logger.Infof("File watcher has started watching files in %s", w.path)
	// We start an infinite loop
	for {
		// `select` statement is like a `switch` but for *channel operations*
		select {
		// Case 1: Read from Events channel
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return // Channel closed, exit goroutine
			}
			w.bufferEvent(event)
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return // Channel closed, exit goroutine
			}
			// Log error but continue watching
			logger.Errorf("Error: %v", err)
		}
	}
}

func (w *Watcher) bufferEvent(event fsnotify.Event) {
	// Only process markdown files and directories
	// TODO: Also process images and pdfs, but leave for later

	if !w.isMarkdownFile(event.Name) && !w.isDirectory(event.Name) {
		return
	}

	if w.isDirectory(event.Name) {
		w.handleDirectoryEvent(event)
		return
	}

	// Buffer file events for debouncing
	now := time.Now()

	// Get or create file event record
	fe, exists := w.eventBuffer[event.Name]

	if !exists {
		fe = &fileEvent{path: event.Name, lastSeen: now}
		w.eventBuffer[event.Name] = fe
	}

	fe.lastSeen = now

	if event.Op&fsnotify.Create == fsnotify.Create {
		logger.Debugf("✅ CREATE event for %s, %s", event.Name, event.Op.String())
		fe.isNew = true
	}

	if event.Op&fsnotify.Write == fsnotify.Write {
		logger.Debugf("🔍 WRITE event for %s, %s", event.Name, event.Op.String())
		fe.isModified = true
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		logger.Debugf("🔍 REMOVE event for %s, %s", event.Name, event.Op.String())
		fe.isDeleted = true
	}

	if event.Op&fsnotify.Rename == fsnotify.Rename {
		logger.Debugf("🔍 RENAME event for %s, %s", event.Name, event.Op.String())
		fe.isDeleted = true // Treat rename as deletion of old name
	}

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	// Process events after 100ms of quiet
	w.debounceTimer = time.AfterFunc(100*time.Millisecond, w.processBufferedEvents)
}

func (w *Watcher) handleDirectoryEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		logger.Infof("📁 New directory created: %s", event.Name)
		if err := w.fsWatcher.Add(event.Name); err != nil {
			logger.Warnf("⚠️ Failed to watch new directory: %s: %v", event.Name, err)
		}
	}
}

func (w *Watcher) processBufferedEvents() {
	now := time.Now()

	for path, fe := range w.eventBuffer {
		// Skip events that are too old
		if now.Sub(fe.lastSeen) > 5*time.Second {
			delete(w.eventBuffer, path)
			continue
		}

		// Determine the primary action
		if fe.isDeleted {
			logger.Infof("🗑️  File deleted: %s", path)
			// TODO: Send delete event to queue

		} else if fe.isNew && !fe.isModified {
			// File was created but not written to (rare)
			logger.Infof("✅ File created (empty): %s", path)
			// TODO: Send create event to queue

		} else if fe.isNew && fe.isModified {
			// File was created and has content (most "new file" cases)
			logger.Infof("✅ File created: %s", path)
			// TODO: Send create event to queue

		} else if fe.isModified {
			// File was modified (existing file edited)
			logger.Infof("✏️  File modified: %s", path)
			// TODO: Send modify event to queue

		} else {
			logger.Warnf("🤷 Unknown event pattern for: %s", path)
		}

		delete(w.eventBuffer, path)
	}
}

func (w *Watcher) Stop() error {
	// Check if fsWatcher is initialized
	if w.fsWatcher != nil {
		close(w.done)
		return w.fsWatcher.Close()
	}

	return nil
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Get file info
	info, err := os.Stat(event.Name)

	// Handle different event types
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		if err == nil && info.IsDir() {
			logger.Infof("📁 New directory created: %s", event.Name)
			if err := w.fsWatcher.Add(event.Name); err != nil {
				logger.Warnf("⚠️ Failed to watch new directory %s: %v", event.Name, err)
			}
		} else if w.isMarkdownFile(event.Name) {
			logger.Infof("✅ File created: %s", event.Name)
		}
	case event.Op&fsnotify.Write == fsnotify.Write:
		if w.isMarkdownFile(event.Name) {
			logger.Infof("✏️ File modified: %s", event.Name)
		}
	}
}

func (w *Watcher) isMarkdownFile(filename string) bool {
	if filepath.Ext(filename) != ".md" {
		return false
	}

	base := filepath.Base(filename)
	if strings.HasPrefix(base, ".") || strings.HasPrefix(base, "~") {
		return false
	}

	return true
}

func (w *Watcher) isDirectory(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && info.IsDir()
}

// Adds a directory and all its subdirectories to the watcher
// addRecursive adds watches recursively to the given root directory and all its subdirectories.
// It returns an error if the root directory cannot be accessed or if there's an issue adding
// watches to any of the directories in the hierarchy.
// addRecursive recursively adds all directories and subdirectories starting from the given root path to the file system watcher.
// It skips any system directories (those prefixed with a dot '.') except for the root directory itself.
// For each path, it logs success or failure of adding the path to the watcher.
//
// Parameters:
//   - root: The starting directory path to begin recursive watching
//
// Returns:
//   - error: Any error that occurs during directory traversal
func (w *Watcher) addRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			logger.Warnf("⚠️ Error accessing %s: %v", path, err)
			return nil
		}

		// Skip system directories, prefixed with '.'
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != root {
				return filepath.SkipDir
			}
		}

		err = w.fsWatcher.Add(path)

		if err != nil {
			logger.Warnf("⚠️ Failed to watch directory %s: %v", path, err)
		} else {
			logger.Debugf("📁 Added directory to watch: %s", path)
		}

		return nil
	})
}
