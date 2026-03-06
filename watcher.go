package main

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	

	"github.com/fsnotify/fsnotify"
)

// isIgnored checks if a directory should be skipped (e.g., .git, node_modules)
func isIgnored(path string) bool {
	ignoredDirs := []string{".git", "node_modules", "bin", "tmp", ".idea"}
	for _, dir := range ignoredDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

// watchDirectories recursively walks the root folder and adds valid directories to the watcher
func watchDirectories(watcher *fsnotify.Watcher, root string) {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// We only watch directories; fsnotify handles the files inside them
		if d.IsDir() {
			if isIgnored(path) {
				return fs.SkipDir
			}
			err = watcher.Add(path)
			if err != nil {
				slog.Error("Failed to watch directory", "path", path, "error", err)
			} else {
				slog.Debug("Watching directory", "path", path)
			}
		}
		return nil
	})

	if err != nil {
		slog.Error("Error walking directories", "error", err)
	}
}