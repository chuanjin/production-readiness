// Package output formats production-readiness findings for display.
// It supports multiple output formats including JSON and Markdown,
// with options for detailed or summary reporting.
package output

import (
	"encoding/json"
	"fmt"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

// ============================================
// JSON Output
// ============================================

// JSONReport represents the structure of the JSON output
type JSONReport struct {
	Summary  SummaryInfo   `json:"summary"`
	Findings FindingsGroup `json:"findings"`
	Signals  *SignalsInfo  `json:"signals,omitempty"`
}

// SummaryInfo contains the overall score and counts
type SummaryInfo struct {
	Score     int `json:"score"`
	Total     int `json:"total"`
	Passed    int `json:"passed"`
	Triggered int `json:"triggered"`
	High      int `json:"high"`
	Medium    int `json:"medium"`
	Low       int `json:"low"`
}

// FindingsGroup groups findings by severity
type FindingsGroup struct {
	High   []FindingDetail `json:"high,omitempty"`
	Medium []FindingDetail `json:"medium,omitempty"`
	Low    []FindingDetail `json:"low,omitempty"`
	Passed []FindingDetail `json:"passed,omitempty"`
}

// FindingDetail represents a single finding in the JSON output
type FindingDetail struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Why         []string `json:"why_it_matters,omitempty"`
	Confidence  string   `json:"confidence,omitempty"`
}

// SignalsInfo contains detected signals from the repository scan
type SignalsInfo struct {
	BoolSignals      map[string]bool   `json:"bool_signals,omitempty"`
	StringSignals    map[string]string `json:"string_signals,omitempty"`
	IntSignals       map[string]int    `json:"int_signals,omitempty"`
	FilesScanned     int               `json:"files_scanned"`
	FilesWithContent int               `json:"files_with_content"`
}

// JSON generates a JSON-formatted report
func JSON(summary engine.Summary, findings []engine.Finding, signals *scanner.RepoSignals) (string, error) {
	report := JSONReport{
		Summary: SummaryInfo{
			Score:     summary.Score,
			Total:     summary.Total,
			Passed:    summary.Passed,
			Triggered: summary.Triggered,
			High:      summary.High,
			Medium:    summary.Medium,
			Low:       summary.Low,
		},
		Findings: FindingsGroup{},
	}

	// Include signals if provided
	if signals != nil {
		report.Signals = &SignalsInfo{
			BoolSignals:      signals.BoolSignals,
			StringSignals:    signals.StringSignals,
			IntSignals:       signals.IntSignals,
			FilesScanned:     len(signals.Files),
			FilesWithContent: len(signals.FileContent),
		}
	}

	// Group findings by severity
	for i := range findings {
		f := &findings[i]

		finding := FindingDetail{
			ID:          f.Rule.ID,
			Title:       f.Rule.Title,
			Description: f.Rule.Description,
			Severity:    string(f.Rule.Severity),
			Category:    f.Rule.Category,
			Why:         f.Rule.Why,
			Confidence:  f.Rule.Confidence,
		}

		if f.Triggered {
			// Add to triggered findings by severity
			switch f.Rule.Severity {
			case "high":
				report.Findings.High = append(report.Findings.High, finding)
			case "medium":
				report.Findings.Medium = append(report.Findings.Medium, finding)
			case "low":
				report.Findings.Low = append(report.Findings.Low, finding)
			}
		} else {
			// Add to passed findings
			report.Findings.Passed = append(report.Findings.Passed, finding)
		}
	}

	// Marshal to JSON with indentation for readability
	jsonBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// JSONCompact generates a compact JSON report (no signals, no passed rules)
func JSONCompact(summary engine.Summary, findings []engine.Finding) (string, error) {
	return JSON(summary, findings, nil)
}
