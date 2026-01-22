package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandleDir(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		root           string
		dirName        string
		ignorePatterns []string
		expectedError  bool // checks if it returns filepath.SkipDir (which is an error value)
	}{
		{
			name:           "Normal directory",
			path:           "/app/src",
			root:           "/app",
			dirName:        "src",
			ignorePatterns: nil,
			expectedError:  false,
		},
		{
			name:           "Default ignored directory (node_modules)",
			path:           "/app/node_modules",
			root:           "/app",
			dirName:        "node_modules",
			ignorePatterns: nil,
			expectedError:  true,
		},
		{
			name:           "Ignored by pattern",
			path:           "/app/build",
			root:           "/app",
			dirName:        "build",
			ignorePatterns: []string{"build/"},
			expectedError:  true,
		},
		{
			name:           "Hidden directory (.git)",
			path:           "/app/.git",
			root:           "/app",
			dirName:        ".git",
			ignorePatterns: nil,
			expectedError:  true,
		},
	}

	logger := &NoopLogger{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleDir(tt.path, tt.root, tt.dirName, tt.ignorePatterns, false, logger)
			if tt.expectedError && err != filepath.SkipDir {
				t.Errorf("expected filepath.SkipDir, got %v", err)
			}
			if !tt.expectedError && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}

func TestHandleFile(t *testing.T) {
	// Create a temporary directory for file operations
	tmpDir := t.TempDir()

	// Create a dummy text file
	textFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(textFile, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a dummy binary file (simulated by extension)
	binFile := filepath.Join(tmpDir, "image.png")
	if err := os.WriteFile(binFile, []byte("fake image data"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create an ignored file
	ignoredFile := filepath.Join(tmpDir, "ignored.log")
	if err := os.WriteFile(ignoredFile, []byte("log data"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name              string
		fileName          string
		ignorePatterns    []string
		shouldBeInFiles   bool
		shouldHaveContent bool
	}{
		{
			name:              "Normal text file",
			fileName:          "test.txt",
			ignorePatterns:    nil,
			shouldBeInFiles:   true,
			shouldHaveContent: true,
		},
		{
			name:              "Binary file",
			fileName:          "image.png",
			ignorePatterns:    nil,
			shouldBeInFiles:   true,
			shouldHaveContent: false,
		},
		{
			name:              "Ignored file",
			fileName:          "ignored.log",
			ignorePatterns:    []string{"*.log"},
			shouldBeInFiles:   true, // It IS added to Files map, but marked as ignored? -> Checked implementation: added to Files map, then checked for ignore.
			shouldHaveContent: false,
		},
	}

	logger := &NoopLogger{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.fileName)
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}

			signals := &RepoSignals{
				Files:           make(map[string]bool),
				FileContent:     make(map[string]string),
				BoolSignals:     make(map[string]bool),
				StringSignals:   make(map[string]string),
				IntSignals:      make(map[string]int),
				DetectedRegions: make(map[string]bool),
			}

			opts := ScanOptions{
				Debug:  false,
				Logger: logger,
			}

			err = handleFile(path, tmpDir, info, tt.ignorePatterns, opts, signals)
			if err != nil {
				t.Fatalf("handleFile returned error: %v", err)
			}

			if signals.Files[tt.fileName] != tt.shouldBeInFiles {
				t.Errorf("Files[%q] = %v, want %v", tt.fileName, signals.Files[tt.fileName], tt.shouldBeInFiles)
			}

			_, hasContent := signals.FileContent[tt.fileName]
			if hasContent != tt.shouldHaveContent {
				t.Errorf("FileContent[%q] exists = %v, want %v", tt.fileName, hasContent, tt.shouldHaveContent)
			}
		})
	}
}
