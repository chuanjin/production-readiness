package output

import (
	"strings"
	"testing"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

func TestMarkdown(t *testing.T) {
	summary := engine.Summary{
		Score:     85,
		Total:     2,
		Passed:    1,
		Triggered: 1,
		High:      1,
	}

	findings := []engine.Finding{
		{
			Triggered: true,
			Rule: rules.Rule{
				ID:          "TEST-001",
				Title:       "High Severity Issue",
				Description: "Fix this critical issue.",
				Severity:    rules.High,
				Why:         []string{"It is dangerous."},
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
			"detected_feature": true,
		},
		StringSignals: map[string]string{
			"version": "1.2.3",
		},
		IntSignals: map[string]int{
			"count": 42,
		},
		Files:       map[string]bool{"main.go": true},
		FileContent: map[string]string{"main.go": "package main"},
	}

	output := Markdown(summary, findings, signals)

	checks := []string{
		"# Production Readiness Report",
		"**Overall Score: 85 / 100**",
		"High Risk",
		"High Severity Issue",
		"Fix this critical issue.",
		"It is dangerous.",
		"Detected Signals",
		"`detected_feature` | ✅",
		"`version` | `1.2.3`",
		"`count` | 42",
		"Files scanned:** 1",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Markdown output missing: %q", check)
		}
	}

	if strings.Contains(output, "Passed Rule") {
		t.Error("Markdown output should not contain passed rules titles in sections")
	}
}

func TestMarkdownSummary(t *testing.T) {
	summary := engine.Summary{
		Score:     90,
		Total:     5,
		Passed:    5,
		Triggered: 0,
	}

	findings := []engine.Finding{
		{Triggered: false, Rule: rules.Rule{Severity: rules.High}},
	}

	output := MarkdownSummary(summary, findings)

	checks := []string{
		"# Production Readiness Summary",
		"**Score: 90 / 100**",
		"✅ No issues found!",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("MarkdownSummary output missing: %q", check)
		}
	}
}
