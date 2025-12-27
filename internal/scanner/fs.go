package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type RepoSignals struct {
	Files       map[string]bool   // filename → exists
	FileContent map[string]string // full path → content
	BoolSignals map[string]bool   // new signals, e.g., "secrets_provider_detected" → true/false
}

func ScanRepo(root string) (RepoSignals, error) {
	signals := RepoSignals{
		Files:       make(map[string]bool),
		FileContent: make(map[string]string),
		BoolSignals: make(map[string]bool),
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := filepath.Base(path)
		signals.Files[name] = true

		if info.Size() < 200_000 { // small text files
			data, err := os.ReadFile(path)
			if err == nil {
				content := string(data)
				if isText(content) {
					signals.FileContent[path] = content
				}
			}
		}
		return nil
	})

	// Later to scan Terraform/Helm manifests or other configs.
	signals.BoolSignals["secrets_provider_detected"] = false

	// Later to detect imports or logging frameworks (e.g., logrus, zap, opentelemetry).
	signals.BoolSignals["correlation_id_detected"] = false
	signals.BoolSignals["structured_logging_detected"] = false

	return signals, err
}

func isText(s string) bool {
	return strings.IndexByte(s, 0) == -1
}
