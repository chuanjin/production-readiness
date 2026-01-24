package test

import (
	"path/filepath"
	"testing"

	"github.com/chuanjin/production-readiness/internal/rules"
)

func TestRulesIntegration(t *testing.T) {
	// Find the rules directory relative to this test file
	// Assuming running from project root or test dir, but hardcoding relation for simplicity
	// "test" package is in <root>/test, "rules" are in <root>/rules

	// Note: When running `go test ./test`, the WD is `.../test`
	rulesDir := "../rules"

	// Attempt to resolve absolute path for clarity in failure messages
	absRulesDir, err := filepath.Abs(rulesDir)
	if err != nil {
		t.Fatalf("Failed to resolve absolute path for rules dir: %v", err)
	}

	loadedRules, err := rules.LoadRules(absRulesDir)
	if err != nil {
		t.Fatalf("Failed to load rules from %s: %v", absRulesDir, err)
	}

	if len(loadedRules) == 0 {
		t.Errorf("No rules loaded from %s. Expected at least one.", absRulesDir)
	}

	// Iterate and verify basic integrity of each rule
	for _, r := range loadedRules {
		if r.ID == "" {
			t.Errorf("Rule found with empty ID")
		}
		if r.Severity == "" {
			t.Errorf("Rule %s has empty Severity", r.ID)
		}
		// Add more checks as needed
	}

	t.Logf("Successfully verified %d rules from %s", len(loadedRules), absRulesDir)
}
