package scanner

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"golang.org/x/sync/errgroup"
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
func ScanRepo(root string) (*RepoSignals, error) {
	return ScanRepoWithOptions(root, ScanOptions{
		Debug:  false,
		Logger: &NoopLogger{},
	})
}

// ScanRepoWithOptions scans the repository with custom options
func ScanRepoWithOptions(root string, opts ScanOptions) (*RepoSignals, error) {
	signals := &RepoSignals{
		Files:           make(map[string]bool),
		FileContent:     make(map[string]string),
		BoolSignals:     make(map[string]bool),
		StringSignals:   make(map[string]string),
		IntSignals:      make(map[string]int),
		DetectedRegions: make(map[string]bool),
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

	// Work channel for parallel processing
	type scanWork struct {
		path  string
		entry os.DirEntry
	}
	workChan := make(chan scanWork, 100)

	g, ctx := errgroup.WithContext(context.Background())

	// Start worker pool
	numWorkers := runtime.NumCPU() * 2
	if numWorkers < 4 {
		numWorkers = 4
	}

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case work, ok := <-workChan:
					if !ok {
						return nil
					}
					if err := handleFile(work.path, root, work.entry, ignorePatterns, opts, signals); err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	// Walk the filesystem and send work to workers
	g.Go(func() error {
		defer close(workChan)
		return filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			// Validate path is within root to prevent traversal
			if !strings.HasPrefix(path, root) {
				return filepath.SkipDir
			}

			// Handle directories
			if info.IsDir() {
				return handleDir(path, root, info.Name(), ignorePatterns, opts.Debug, logger)
			}

			// Send file to work channel
			select {
			case workChan <- scanWork{path: path, entry: info}:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	err := g.Wait()

	if opts.Debug {
		logger.Println("\n=== Summary ===")
		// Using helper methods since maps are now protected by mutex
		logger.Printf("Total files: %d", len(signals.GetFiles()))
		logger.Printf("Files with content: %d", len(signals.GetFileContentMap()))
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
