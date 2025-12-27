package engine

import (
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
	Score    int
	Findings []Finding
}

// ConditionFunc evaluates a single condition
type ConditionFunc func(value interface{}, signals scanner.RepoSignals) bool

// ConditionRegistry holds all available condition handlers
var ConditionRegistry = map[string]ConditionFunc{
	"file_exists": func(value interface{}, signals scanner.RepoSignals) bool {
		name := value.(string)
		return signals.Files[name]
	},
	"code_contains": func(value interface{}, signals scanner.RepoSignals) bool {
		needle := value.(string)
		for _, content := range signals.FileContent {
			if strings.Contains(content, needle) {
				return true
			}
		}
		return false
	},
	"secrets_provider_detected": func(value interface{}, signals scanner.RepoSignals) bool {
		expected := value.(bool)
		return signals.BoolSignals["secrets_provider_detected"] == expected
	},
	"correlation_id_detected": func(value interface{}, signals scanner.RepoSignals) bool {
		expected := value.(bool)
		return signals.BoolSignals["correlation_id_detected"] == expected
	},
	"structured_logging_detected": func(value interface{}, signals scanner.RepoSignals) bool {
		expected := value.(bool)
		return signals.BoolSignals["structured_logging_detected"] == expected
	},
	"infra_as_code_detected": func(value interface{}, signals scanner.RepoSignals) bool {
		expected := value.(bool)
		return signals.BoolSignals["infra_as_code_detected"] == expected
	},
	// Add string signal check
	"some_string_signal": func(value interface{}, signals scanner.RepoSignals) bool {
		expected := value.(string)
		return signals.StringSignals["some_string_signal"] == expected
	},
}

// Evaluate all rules and return a summary
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

func Summarize(findings []Finding) Summary {
	summary := Summary{Findings: findings}

	for _, f := range findings {
		if !f.Supported {
			continue
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

	summary.Score = 100 - (summary.High*20 + summary.Medium*10 + summary.Low*5)
	if summary.Score < 0 {
		summary.Score = 0
	}
	return summary
}

func evaluateRule(rule rules.Rule, signals scanner.RepoSignals) (triggered bool, supported bool) {
	supported = false
	triggered = false

	// Determine if rule is supported by checking if all its condition keys exist in registry
	allKeys := append(append(rule.Detect.AnyOf, rule.Detect.AllOf...), rule.Detect.NoneOf...)
	for _, cond := range allKeys {
		for k := range cond {
			if _, ok := ConditionRegistry[k]; ok {
				supported = true
			}
		}
	}

	// Evaluate AnyOf
	for _, cond := range rule.Detect.AnyOf {
		if matchCondition(cond, signals) {
			triggered = true
			break
		}
	}

	// Evaluate AllOf
	for _, cond := range rule.Detect.AllOf {
		if !matchCondition(cond, signals) {
			triggered = false
			break
		}
	}

	// Evaluate NoneOf
	for _, cond := range rule.Detect.NoneOf {
		if matchCondition(cond, signals) {
			triggered = false
			break
		}
	}

	return triggered, supported
}

func matchCondition(cond map[string]interface{}, signals scanner.RepoSignals) bool {
	for key, val := range cond {
		if fn, ok := ConditionRegistry[key]; ok {
			if fn(val, signals) {
				return true
			}
		} else {
			// Unknown condition: treat as unsupported
			continue
		}
	}
	return false
}
