package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// handleDir processes directories during walk
func handleDir(path, root, name string, ignorePatterns []string, debug bool, logger Logger) error {
	if defaultIgnoredDirs[name] {
		return filepath.SkipDir
	}

	rel, e := filepath.Rel(root, path)
	if e != nil {
		return filepath.SkipDir
	}

	if isIgnored(rel, ignorePatterns) {
		if debug {
			logger.Printf("Skipping directory: %s (ignored)", rel)
		}
		return filepath.SkipDir
	}
	return nil
}

// handleFile processes files during walk
func handleFile(path, root string, info os.DirEntry, ignorePatterns []string, opts ScanOptions, signals *RepoSignals) error {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return nil
	}

	if opts.Debug {
		opts.Logger.Printf("Processing: %s", relPath)
	}

	// Always store to Files for existing check in engine
	signals.SetFile(relPath)
	if opts.Debug {
		opts.Logger.Println("  -> Added to Files map")
	}

	// Check if ignored BEFORE adding to Files
	if isIgnored(relPath, ignorePatterns) {
		if opts.Debug {
			opts.Logger.Println("  -> IGNORED by .prignore")
		}
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))

	if binaryExts[ext] {
		if opts.Debug {
			opts.Logger.Println("  -> Skipped (binary)")
		}
		return nil
	}

	finfo, err := info.Info()
	if err != nil {
		return nil
	}

	if finfo.Size() < 200_000 {
		// #nosec G304 - path is validated to be within root directory
		data, err := os.ReadFile(path)
		if err == nil && isText(string(data)) {
			content := string(data)
			signals.SetContent(relPath, content)
			if opts.Debug {
				opts.Logger.Println("  -> Added to FileContent")
			}

			// Run all detectors dynamically
			runAllDetectors(content, relPath, signals)
		}
	} else if opts.Debug {
		opts.Logger.Println("  -> Skipped (too large)")
	}
	return nil
}
