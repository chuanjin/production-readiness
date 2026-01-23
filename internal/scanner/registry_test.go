package scanner

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	// Verify that init() has registered some detectors
	if len(detectorRegistry) == 0 {
		t.Error("Expected initial detectors to be registered, but registry is empty")
	}

	// Test adding a custom detector and running it
	var called bool
	var capturedContent, capturedPath string

	mockDetector := func(content string, relPath string, signals *RepoSignals) {
		called = true
		capturedContent = content
		capturedPath = relPath
	}

	// Register the mock detector
	registerDetector(mockDetector)

	// Run all detectors
	testContent := "test content"
	testPath := "test.txt"
	signals := &RepoSignals{
		Files:           make(map[string]bool),
		FileContent:     make(map[string]string),
		BoolSignals:     make(map[string]bool),
		StringSignals:   make(map[string]string),
		IntSignals:      make(map[string]int),
		DetectedRegions: make(map[string]bool),
	}

	runAllDetectors(testContent, testPath, signals)

	if !called {
		t.Error("Expected mock detector to be called, but it wasn't")
	}
	if capturedContent != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, capturedContent)
	}
	if capturedPath != testPath {
		t.Errorf("Expected path '%s', got '%s'", testPath, capturedPath)
	}
}
