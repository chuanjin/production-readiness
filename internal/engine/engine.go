package engine

import (
	"path/filepath"
	"strings"

	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

type Finding struct {
	Rule      rules.Rule
	Triggered bool
	Supported bool
}

type Summary struct {
	High     int
	Medium   int
	Low      int
	Positive int

	Findings []Finding
	Skipped  []Finding
	NotHit   []Finding

	Score int
}

func Evaluate(ruleset []rules.Rule, signals scanner.RepoSignals) Summary {
	var findings []Finding

	for _, r := range ruleset {
		triggered, supported := evaluateRule(r, signals)
		f := Finding{
			Rule:      r,
			Triggered: triggered,
			Supported: supported,
		}
		findings = append(findings, f)
	}

	return Summarize(findings)
}

// RETURNS (triggered, supported)
func evaluateRule(rule rules.Rule, signals scanner.RepoSignals) (bool, bool) {
	supported := false

	conds := rule.Detect.AnyOf
	if len(conds) == 0 {
		return false, false
	}

	for _, cond := range conds {
		// supported types
		if cond["file_exists"] != nil || cond["code_contains"] != nil {
			supported = true
		}

		// file_exists
		if val, ok := cond["file_exists"]; ok {
			expect := filepath.Clean(val.(string))

			// check basename (example: "Dockerfile")
			if _, exists := signals.Files[expect]; exists {
				return true, true
			}

			// wildcard (example "*.yaml")
			for file := range signals.Files {
				match, _ := filepath.Match(expect, file)
				if match {
					return true, true
				}
			}
		}

		// code_contains
		if val, ok := cond["code_contains"]; ok {
			needle := val.(string)

			for _, content := range signals.FileContent {
				if strings.Contains(content, needle) {
					return true, true
				}
			}
		}
	}

	return false, supported
}

func Summarize(findings []Finding) Summary {
	sum := Summary{}
	sum.Findings = findings

	for _, f := range findings {
		if !f.Supported {
			sum.Skipped = append(sum.Skipped, f)
			continue
		}

		if !f.Triggered {
			sum.NotHit = append(sum.NotHit, f)
			continue
		}

		switch f.Rule.Severity {
		case "high":
			sum.High++
		case "medium":
			sum.Medium++
		case "low":
			sum.Low++
		default:
			sum.Positive++
		}
	}

	sum.Score = 100 - (sum.High*20 + sum.Medium*10 + sum.Low*5)
	if sum.Score < 0 {
		sum.Score = 0
	}

	return sum
}
