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
	Supported bool // true if we know how to detect this rule
}

type Summary struct {
	High     int
	Medium   int
	Low      int
	Positive int
	Score    int
	Findings []Finding
}

// Evaluate all rules against repo signals
func Evaluate(ruleSet []rules.Rule, signals scanner.RepoSignals) Summary {
	var findings []Finding

	for _, r := range ruleSet {
		triggered, supported := evaluateRule(r, signals)
		findings = append(findings, Finding{
			Rule:      r,
			Triggered: triggered,
			Supported: supported,
		})
	}

	return Summarize(findings)
}

// Summarize converts findings to a summary struct with score
func Summarize(findings []Finding) Summary {
	summary := Summary{
		Findings: findings,
	}

	for _, f := range findings {
		if !f.Supported {
			continue // skip unsupported rules
		}
		switch f.Rule.Severity {
		case rules.High:
			if f.Triggered {
				summary.High++
			}
		case rules.Medium:
			if f.Triggered {
				summary.Medium++
			}
		case rules.Low:
			if f.Triggered {
				summary.Low++
			}
		case rules.Positive:
			if f.Triggered {
				summary.Positive++
			}
		}
	}

	// simple scoring formula
	summary.Score = 100 - (summary.High*20 + summary.Medium*10 + summary.Low*5)
	if summary.Score < 0 {
		summary.Score = 0
	}

	return summary
}

// evaluateRule determines if a rule triggers, and if it is supported
func evaluateRule(rule rules.Rule, signals scanner.RepoSignals) (triggered bool, supported bool) {
	supported = false
	triggered = false

	// Determine if rule has any supported conditions
	for _, cond := range rule.Detect.AnyOf {
		if cond["file_exists"] != nil || cond["code_contains"] != nil || cond["secrets_provider_detected"] != nil ||
			cond["correlation_id_detected"] != nil || cond["structured_logging_detected"] != nil {
			supported = true
		}
	}
	for _, cond := range rule.Detect.AllOf {
		if cond["file_exists"] != nil || cond["code_contains"] != nil || cond["secrets_provider_detected"] != nil ||
			cond["correlation_id_detected"] != nil || cond["structured_logging_detected"] != nil {
			supported = true
		}
	}
	for _, cond := range rule.Detect.NoneOf {
		if cond["file_exists"] != nil || cond["code_contains"] != nil || cond["secrets_provider_detected"] != nil ||
			cond["correlation_id_detected"] != nil || cond["structured_logging_detected"] != nil {
			supported = true
		}
	}

	// Evaluate any_of: triggers if any condition matches
	for _, cond := range rule.Detect.AnyOf {
		if matchCondition(cond, signals) {
			triggered = true
			break
		}
	}

	// Evaluate all_of: must match all conditions
	for _, cond := range rule.Detect.AllOf {
		if !matchCondition(cond, signals) {
			triggered = false
			break
		}
	}

	// Evaluate none_of: un-trigger if any none_of condition matches
	for _, cond := range rule.Detect.NoneOf {
		if matchCondition(cond, signals) {
			triggered = false
			break
		}
	}

	return triggered, supported
}

// matchCondition evaluates a single condition against signals
func matchCondition(cond map[string]interface{}, signals scanner.RepoSignals) bool {
	for k, v := range cond {
		switch k {
		case "file_exists":
			name := v.(string)
			if signals.Files[name] {
				return true
			}
			for file := range signals.Files {
				match, _ := filepath.Match(name, file)
				if match {
					return true
				}
			}

		case "code_contains":
			needle := v.(string)
			for _, content := range signals.FileContent {
				if strings.Contains(content, needle) {
					return true
				}
			}

		case "secrets_provider_detected", "correlation_id_detected", "structured_logging_detected":
			expect := v.(bool)
			if val, ok := signals.BoolSignals[k]; ok && val == expect {
				return true
			}
		}
	}

	return false
}
