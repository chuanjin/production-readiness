package output

import (
	"encoding/json"

	"github.com/chuanjin/production-readiness/internal/engine"
)

// JSONReport represents the structure of the JSON output
type JSONReport struct {
	Summary  SummaryInfo   `json:"summary"`
	Findings FindingsGroup `json:"findings"`
}

// SummaryInfo contains the overall score and counts
type SummaryInfo struct {
	Score  int `json:"score"`
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

// FindingsGroup groups findings by severity
type FindingsGroup struct {
	High   []FindingDetail `json:"high,omitempty"`
	Medium []FindingDetail `json:"medium,omitempty"`
	Low    []FindingDetail `json:"low,omitempty"`
}

// FindingDetail represents a single finding in the JSON output
type FindingDetail struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category,omitempty"`
	Why         []string `json:"why_it_matters,omitempty"`
	Confidence  string   `json:"confidence,omitempty"`
}

// JSON generates a JSON-formatted report
func JSON(summary engine.Summary, findings []engine.Finding) string {
	report := JSONReport{
		Summary: SummaryInfo{
			Score:  summary.Score,
			High:   summary.High,
			Medium: summary.Medium,
			Low:    summary.Low,
		},
		Findings: FindingsGroup{},
	}

	// Group findings by severity (only triggered ones)
	for _, f := range findings {
		if !f.Triggered {
			continue
		}

		finding := FindingDetail{
			ID:          f.Rule.ID,
			Title:       f.Rule.Title,
			Description: f.Rule.Description,
			Category:    f.Rule.Category,
			Why:         f.Rule.Why,
			Confidence:  f.Rule.Confidence,
		}

		switch f.Rule.Severity {
		case "high":
			report.Findings.High = append(report.Findings.High, finding)
		case "medium":
			report.Findings.Medium = append(report.Findings.Medium, finding)
		case "low":
			report.Findings.Low = append(report.Findings.Low, finding)
		}
	}

	// Marshal to JSON with indentation for readability
	jsonBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		// Fallback to compact JSON if indentation fails
		jsonBytes, _ = json.Marshal(report)
	}

	return string(jsonBytes)
}
