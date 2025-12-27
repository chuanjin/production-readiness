package scanner

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RepoSignals struct {
	Files         map[string]bool
	FileContent   map[string]string
	BoolSignals   map[string]bool
	StringSignals map[string]string
}

func ScanRepo(root string) (RepoSignals, error) {
	signals := RepoSignals{
		Files:         make(map[string]bool),
		FileContent:   make(map[string]string),
		BoolSignals:   make(map[string]bool),
		StringSignals: make(map[string]string),
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("error accessing %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
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

	return signals, err
}

func isText(s string) bool {
	return strings.IndexByte(s, 0) == -1
}
