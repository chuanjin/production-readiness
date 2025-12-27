package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadPrIgnore reads the .prignore file in the repo root.
// Returns a list of skip patterns (empty if no file).
func LoadPrIgnore(root string) ([]string, error) {
	var patterns []string
	path := filepath.Join(root, ".prignore")
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return patterns, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// remove inline comments after space
		if idx := strings.Index(line, " "); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		patterns = append(patterns, line)
	}
	return patterns, scanner.Err()
}

// ShouldIgnore checks if a file path matches any pattern in skipList
func ShouldIgnore(path string, skipList []string, root string) bool {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		relPath = path
	}
	relPath = filepath.ToSlash(relPath) // normalize slashes

	for _, pattern := range skipList {
		pattern = filepath.ToSlash(pattern)
		// folder pattern
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(relPath, strings.TrimSuffix(pattern, "/")+"/") {
				return true
			}
		} else {
			// glob pattern match
			match, err := filepath.Match(pattern, filepath.Base(relPath))
			if err == nil && match {
				return true
			}
		}
	}
	return false
}
