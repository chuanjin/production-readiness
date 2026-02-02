package scanner

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsText(t *testing.T) {
	if !isText("hello world") {
		t.Fatalf("expected plain text to be detected as text")
	}
	if isText("binary\x00data") {
		t.Fatalf("expected string with NUL byte to be detected as non-text")
	}
}

func TestParsePrIgnore(t *testing.T) {
	dir := t.TempDir()
	content := "# comment line\n\nignoreme/**\n*.md\nlogs/\n"
	if err := os.WriteFile(filepath.Join(dir, ".prignore"), []byte(content), 0o644); err != nil {
		t.Fatalf("write .prignore: %v", err)
	}

	got := parsePrIgnore(dir)
	want := []string{"ignoreme/**", "*.md", "logs/"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestIsIgnored(t *testing.T) {
	patterns := []string{"*.yaml", "docs/**", "build/", "secret.txt"}

	cases := []struct {
		path     string
		expected bool
	}{
		{path: "config.yaml", expected: true},              // basename match
		{path: "rules/00-example.yaml", expected: true},    // basename match with subdir
		{path: "docs/setup/guide.md", expected: true},      // nested directory match
		{path: "build/output.bin", expected: true},         // trailing slash directory match
		{path: "secret.txt", expected: true},               // direct filename match
		{path: "src/main.go", expected: false},             // not ignored
		{path: "documentation/readme.md", expected: false}, // prefix only, not matching pattern
	}

	for _, tt := range cases {
		if isIgnored(tt.path, patterns) != tt.expected {
			t.Fatalf("isIgnored(%q) expected %v", tt.path, tt.expected)
		}
	}
}

func TestScanRepo(t *testing.T) {
	root := t.TempDir()

	// .prignore to skip markdown and a specific directory
	prignore := "ignoreme/**\n*.md\n"
	if err := os.WriteFile(filepath.Join(root, ".prignore"), []byte(prignore), 0o644); err != nil {
		t.Fatalf("write .prignore: %v", err)
	}

	// File that should be scanned and parsed (also triggers a detector)
	mainContent := `region = "us-east-1"`
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(mainContent), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	// Ignored by pattern but still recorded in Files map
	if err := os.Mkdir(filepath.Join(root, "ignoreme"), 0o755); err != nil {
		t.Fatalf("mkdir ignoreme: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignoreme", "ignored.txt"), []byte("should ignore"), 0o644); err != nil {
		t.Fatalf("write ignored.txt: %v", err)
	}

	// Markdown ignored by .prignore
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("write README.md: %v", err)
	}

	// Binary file should be skipped from FileContent
	if err := os.WriteFile(filepath.Join(root, "image.png"), []byte{0x89, 0x50, 0x4E, 0x47}, 0o644); err != nil {
		t.Fatalf("write image.png: %v", err)
	}

	// Large file should be skipped due to size
	largeData := bytes.Repeat([]byte("a"), 210_000)
	if err := os.WriteFile(filepath.Join(root, "large.txt"), largeData, 0o644); err != nil {
		t.Fatalf("write large.txt: %v", err)
	}

	// Default ignored directory (node_modules) should be skipped entirely
	if err := os.Mkdir(filepath.Join(root, "node_modules"), 0o755); err != nil {
		t.Fatalf("mkdir node_modules: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "node_modules", "dep.js"), []byte("console.log('hi')"), 0o644); err != nil {
		t.Fatalf("write dep.js: %v", err)
	}

	signals, err := ScanRepo(root)
	if err != nil {
		t.Fatalf("ScanRepo returned error: %v", err)
	}

	// Files map should include all files except those under default ignored directories
	if !signals.Files["main.go"] {
		t.Fatalf("expected main.go to be present in Files map")
	}
	if !signals.Files["README.md"] {
		t.Fatalf("expected README.md to be present in Files map even if ignored")
	}
	if signals.Files["node_modules/dep.js"] {
		t.Fatalf("expected node_modules to be skipped entirely")
	}

	// FileContent should only include scanned text files not ignored by .prignore, not binary, not oversized
	if _, ok := signals.FileContent["main.go"]; !ok {
		t.Fatalf("expected main.go content to be captured")
	}
	if _, ok := signals.FileContent["README.md"]; ok {
		t.Fatalf("expected README.md to be ignored by .prignore")
	}
	if _, ok := signals.FileContent["ignoreme/ignored.txt"]; ok {
		t.Fatalf("expected ignoreme/ignored.txt to be ignored by .prignore")
	}
	if _, ok := signals.FileContent["image.png"]; ok {
		t.Fatalf("expected binary file to be skipped")
	}
	if _, ok := signals.FileContent["large.txt"]; ok {
		t.Fatalf("expected large file to be skipped")
	}

	// Detectors should have run on scanned content
	if got := signals.GetInt("region_count"); got != 1 {
		t.Fatalf("expected region_count to be 1, got %d", got)
	}
}
