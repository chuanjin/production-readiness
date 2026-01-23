package version

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	// Verify that the version variables are accessible and have default values
	// Since these are set at build time, we expect the defaults here in tests
	if Version != "dev" {
		t.Errorf("Expected Version to be 'dev', got '%s'", Version)
	}

	if Commit != "none" {
		t.Errorf("Expected Commit to be 'none', got '%s'", Commit)
	}

	if BuildDate != "unknown" {
		t.Errorf("Expected BuildDate to be 'unknown', got '%s'", BuildDate)
	}

	// Verify we can modify them (simulating build flags)
	originalVersion := Version
	defer func() { Version = originalVersion }()

	Version = "1.0.0"
	if Version != "1.0.0" {
		t.Errorf("Expected Version to be updated to '1.0.0', got '%s'", Version)
	}
}
