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

// Logger interface allows for flexible logging implementations
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// ScanOptions configures the repository scan
type ScanOptions struct {
	Debug  bool
	Logger Logger
}

// NoopLogger is a no-op logger (exported for external use)
type NoopLogger struct{}

func (n *NoopLogger) Printf(format string, v ...interface{}) {}
func (n *NoopLogger) Println(v ...interface{})               {}

// ScanRepo scans the repository with default options (no debug output)
func ScanRepo(root string) (RepoSignals, error) {
	return ScanRepoWithOptions(root, ScanOptions{
		Debug:  false,
		Logger: &NoopLogger{},
	})
}

// ScanRepoWithOptions scans the repository with custom options
func ScanRepoWithOptions(root string, opts ScanOptions) (RepoSignals, error) {
	signals := RepoSignals{
		Files:         make(map[string]bool),
		FileContent:   make(map[string]string),
		BoolSignals:   make(map[string]bool),
		StringSignals: make(map[string]string),
		IntSignals:    make(map[string]int),
	}

	// Use provided logger or default to noop
	logger := opts.Logger
	if logger == nil {
		logger = &NoopLogger{}
	}

	ignorePatterns := parsePrIgnore(root)

	if opts.Debug {
		logger.Println("=== Ignore patterns ===")
		for _, p := range ignorePatterns {
			logger.Printf("Pattern: %s", p)
		}
		logger.Println("")
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Validate path is within root to prevent traversal
		if !strings.HasPrefix(path, root) {
			return filepath.SkipDir
		}

		// Skip ignored folders before entering them
		if info.IsDir() {
			name := info.Name()
			if defaultIgnoredDirs[name] {
				return filepath.SkipDir
			}

			rel, e := filepath.Rel(root, path)
			if e != nil {
				return filepath.SkipDir
			}

			if isIgnored(rel, ignorePatterns) {
				if opts.Debug {
					logger.Printf("Skipping directory: %s (ignored)", rel)
				}
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		if opts.Debug {
			logger.Printf("Processing: %s", relPath)
		}

		// Always store to Files for existing check in engine
		signals.Files[relPath] = true
		if opts.Debug {
			logger.Println("  -> Added to Files map")
		}

		// Check if ignored BEFORE adding to Files
		if isIgnored(relPath, ignorePatterns) {
			if opts.Debug {
				logger.Println("  -> IGNORED by .prignore")
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if binaryExts[ext] {
			if opts.Debug {
				logger.Println("  -> Skipped (binary)")
			}
			return nil
		}

		if info.Size() < 200_000 {
			// #nosec G304 - path is validated to be within root directory
			data, err := os.ReadFile(path)
			if err == nil && isText(string(data)) {
				content := string(data)
				signals.FileContent[relPath] = content
				if opts.Debug {
					logger.Println("  -> Added to FileContent")
				}

				// Run all detectors dynamically
				runAllDetectors(content, relPath, &signals)
			}
		} else if opts.Debug {
			logger.Println("  -> Skipped (too large)")
		}
		return nil
	})

	if opts.Debug {
		logger.Println("\n=== Summary ===")
		logger.Printf("Total files: %d", len(signals.Files))
		logger.Printf("Files with content: %d", len(signals.FileContent))
	}

	return signals, err
}

func isText(s string) bool {
	return strings.IndexByte(s, 0) == -1
}

// parsePrIgnore reads .prignore and converts lines to matchable patterns
func parsePrIgnore(root string) []string {
	var ignores []string

	// Validate root is clean
	root, err := filepath.Abs(root)
	if err != nil {
		return ignores
	}

	path := filepath.Join(root, ".prignore")

	// Ensure path is still within root
	if !strings.HasPrefix(path, root) {
		return ignores
	}

	// #nosec G304 - path is validated to be within root directory
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
		match, err := doublestar.Match(pattern, relPath)
		if err != nil {
			continue
		}
		if match {
			return true
		}

		// Also try matching against just the basename
		// This allows "*.yaml" to match "rules/00-example.yaml"
		match, err = doublestar.Match(pattern, basename)
		if err != nil {
			continue
		}
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
