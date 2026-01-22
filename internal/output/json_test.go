package output

import (
	"encoding/json"
	"testing"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

func TestJSON(t *testing.T) {
	summary := engine.Summary{
		Score:     80,
		Total:     2,
		Passed:    1,
		Triggered: 1,
		High:      1,
	}

	findings := []engine.Finding{
		{
			Triggered: true,
			Rule: rules.Rule{
				ID:       "TEST-001",
				Title:    "High Severity Issue",
				Severity: rules.High,
			},
		},
		{
			Triggered: false,
			Rule: rules.Rule{
				ID:       "TEST-002",
				Title:    "Passed Rule",
				Severity: rules.Medium,
			},
		},
	}

	signals := &scanner.RepoSignals{
		BoolSignals: map[string]bool{
			"test_signal": true,
		},
		Files:       map[string]bool{"main.go": true},
		FileContent: map[string]string{"main.go": "package main"},
	}

	t.Run("Full Report", func(t *testing.T) {
		output, err := JSON(summary, findings, signals)
		if err != nil {
			t.Fatalf("JSON() error = %v", err)
		}

		// Verify structure using unmarshal
		var report JSONReport
		if err := json.Unmarshal([]byte(output), &report); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v", err)
		}

		if report.Summary.Score != 80 {
			t.Errorf("Expected score 80, got %d", report.Summary.Score)
		}

		if len(report.Findings.High) != 1 {
			t.Errorf("Expected 1 high finding, got %d", len(report.Findings.High))
		}

		if len(report.Findings.Passed) != 1 {
			t.Errorf("Expected 1 passed finding, got %d", len(report.Findings.Passed))
		}

		if report.Signals == nil || !report.Signals.BoolSignals["test_signal"] {
			t.Error("Signals not correctly included")
		}
	})

	t.Run("Compact Report", func(t *testing.T) {
		output, err := JSONCompact(summary, findings)
		if err != nil {
			t.Fatalf("JSONCompact() error = %v", err)
		}

		var report JSONReport
		if err := json.Unmarshal([]byte(output), &report); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v", err)
		}

		if report.Signals != nil {
			t.Error("Compact report should have nil Signals")
		}

		if len(report.Findings.Passed) > 0 {
			t.Errorf("Compact report should not contain passed findings, got %d", len(report.Findings.Passed))
		}
	})
}
