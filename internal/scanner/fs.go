package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// binaryExts are automatically skipped
var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".pdf": true, ".zip": true, ".tar": true, ".gz": true,
	".mp4": true, ".mp3": true, ".mov": true, ".exe": true,
}

// defaultIgnoredDirs — skipped always (low-level ignore, even if .prignore missing)
var defaultIgnoredDirs = map[string]bool{
	".git":         true,
	".svn":         true,
	".hg":          true,
	"node_modules": true,
}

func ScanRepo(root string) (RepoSignals, error) {
	signals := RepoSignals{
		Files:         make(map[string]bool),
		FileContent:   make(map[string]string),
		BoolSignals:   make(map[string]bool),
		StringSignals: make(map[string]string),
		IntSignals:    make(map[string]int),
	}

	ignorePatterns := parsePrIgnore(root)

	println("=== Ignore patterns ===")
	for _, p := range ignorePatterns {
		println("Pattern:", p)
	}
	println("")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip ignored folders before entering them
		if info.IsDir() {
			name := info.Name()
			if defaultIgnoredDirs[name] {
				return filepath.SkipDir
			}
			rel, _ := filepath.Rel(root, path)
			if isIgnored(rel, ignorePatterns) {
				return filepath.SkipDir
			}
			return nil // continue walking
		}

		relPath, _ := filepath.Rel(root, path)

		println("Processing:", relPath)

		// Always store to Files for existing check in engine
		signals.Files[relPath] = true
		println("  -> Added to Files map")

		// Check if ignored BEFORE adding to Files
		if isIgnored(relPath, ignorePatterns) {
			println("  -> IGNORED by .prignore")
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if binaryExts[ext] {
			println("  -> Skipped (binary)")
			return nil
		}

		if info.Size() < 200_000 {
			data, err := os.ReadFile(path)
			if err == nil && isText(string(data)) {
				content := string(data)
				signals.FileContent[relPath] = content
				println("  -> Added to FileContent")

				// Run all detectors dynamically
				runAllDetectors(content, relPath, &signals)
			}
		} else {
			println("  -> Skipped (too large)")
		}
		return nil
	})

	println("\n=== Summary ===")
	println("Total files:", len(signals.Files))
	println("Files with content:", len(signals.FileContent))

	return signals, err
}

func isText(s string) bool {
	return strings.IndexByte(s, 0) == -1
}

// parsePrIgnore reads .prignore and converts lines to matchable patterns
func parsePrIgnore(root string) []string {
	var ignores []string
	path := filepath.Join(root, ".prignore")
	data, err := os.ReadFile(path)
	if err != nil {
		return ignores
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Normalize separators for cross-platform
		line = filepath.ToSlash(line)
		ignores = append(ignores, line)
	}

	return ignores
}

// isIgnored checks if a relative path matches any .prignore pattern
func isIgnored(relPath string, ignorePatterns []string) bool {
	relPath = filepath.ToSlash(relPath)
	basename := filepath.Base(relPath)

	for _, pattern := range ignorePatterns {
		// Try matching against full path
		match, _ := doublestar.Match(pattern, relPath)
		if match {
			return true
		}

		// Also try matching against just the basename
		// This allows "*.yaml" to match "rules/00-example.yaml"
		match, _ = doublestar.Match(pattern, basename)
		if match {
			return true
		}

		// Handle directory patterns (those ending with /)
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, dirPattern+"/") || relPath == dirPattern {
				return true
			}
		}
	}
	return false
}
