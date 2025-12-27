package scanner

import (
	"os"
	"path/filepath"
	"strings"
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

// ScanRepo scans the repository at root, respecting skipList (.prignore)
// Files in skipList still populate FileExists, but content is skipped
func ScanRepo(root string, skipList []string) (RepoSignals, error) {
	signals := RepoSignals{
		Files:         make(map[string]bool),
		FileContent:   make(map[string]string),
		BoolSignals:   make(map[string]bool),
		StringSignals: make(map[string]string),
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := filepath.Base(path)
		signals.Files[name] = true

		// skip reading content if ignored
		if ShouldIgnore(path, skipList, root) {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if binaryExts[ext] {
			return nil
		}

		if !allowedCodeExts[ext] {
			return nil
		}

		if info.Size() < 200_000 { // only read small files
			data, err := os.ReadFile(path)
			if err == nil && isText(string(data)) {
				signals.FileContent[path] = string(data)
			}
		}
		return nil
	})

	return signals, err
}

// isText checks for null bytes
func isText(s string) bool {
	return strings.IndexByte(s, 0) == -1
}
