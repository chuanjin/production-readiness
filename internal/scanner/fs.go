package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var allowedCodeExts = map[string]bool{
	// existing code files
	".go": true, ".ts": true, ".js": true, ".jsx": true, ".tsx": true,
	".py": true, ".rb": true, ".java": true, ".kt": true, ".swift": true,
	".c": true, ".cpp": true, ".h": true, ".rs": true, ".php": true,
	".sh": true, ".bash": true, ".zsh": true, ".tf": true, ".dockerfile": true,
	".sql": true, ".gradle": true, ".makefile": true,
	"": true, // extensionless like Dockerfile / Makefile

	// Add config files
	".env":  true,               // environment variables
	".yaml": true, ".yml": true, // k8s / helm / config
	".json": true, // package.json, deployment config
	".toml": true, // configs like pyproject.toml
	".ini":  true, // generic config
}

// binaryExts are automatically skipped
var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".pdf": true, ".zip": true, ".tar": true, ".gz": true,
	".mp4": true, ".mp3": true, ".mov": true, ".exe": true,
}

func ScanRepo(root string) (RepoSignals, error) {
	signals := RepoSignals{
		Files:         make(map[string]bool),
		FileContent:   make(map[string]string),
		BoolSignals:   make(map[string]bool),
		StringSignals: make(map[string]string),
		IntSignals:    make(map[string]int),
	}

	ignorePaths := parsePrIgnore(root)

	println("=== Ignore patterns ===")
	for _, p := range ignorePaths {
		println("Pattern:", p)
	}
	println("")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)

		println("Processing:", relPath)

		// Check if ignored BEFORE adding to Files
		if isIgnored(relPath, ignorePaths) {
			println("  -> IGNORED by .prignore")
			return nil
		}

		// Store the relative path
		signals.Files[relPath] = true
		println("  -> Added to Files map")

		ext := strings.ToLower(filepath.Ext(path))

		if binaryExts[ext] {
			println("  -> Skipped (binary)")
			return nil
		}

		if !allowedCodeExts[ext] && !allowedCodeExts[info.Name()] {
			println("  -> Skipped (not allowed extension)")
			return nil
		}

		if info.Size() < 200_000 {
			data, err := os.ReadFile(path)
			if err == nil && isText(string(data)) {
				signals.FileContent[relPath] = string(data)
				println("  -> Added to FileContent")
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
