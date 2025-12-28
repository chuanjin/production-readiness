package engine

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

// ConditionFunc evaluates a condition
type ConditionFunc func(value interface{}, signals scanner.RepoSignals) bool

// ConditionRegistry holds all registered condition functions
var ConditionRegistry = map[string]ConditionFunc{}

func init() {
	// ===== Built-in evaluators =====

	// file_exists: ".env" OR "*.env" OR "**/*.yaml"
	ConditionRegistry["file_exists"] = func(value interface{}, signals scanner.RepoSignals) bool {
		pattern := value.(string)

		// 1️⃣ basename exact match
		for full := range signals.Files {
			if filepath.Base(full) == pattern {
				return true
			}
		}

		// 2️⃣ glob match using doublestar across full repo paths
		for full := range signals.Files {
			match, _ := doublestar.Match(pattern, full)
			if match {
				return true
			}
		}
		return false
	}

	// code_contains: "os.environ"
	ConditionRegistry["code_contains"] = func(value interface{}, signals scanner.RepoSignals) bool {
		needle := value.(string)
		for _, content := range signals.FileContent {
			if strings.Contains(content, needle) {
				return true
			}
		}
		return false
	}

	ConditionRegistry["signal_equals"] = func(value interface{}, signals scanner.RepoSignals) bool {
		params := value.(map[string]interface{})
		for key, expected := range params {

			// bool signal
			if actual, ok := signals.BoolSignals[key]; ok {
				return actual == expected
			}

			// string signal
			if actual, ok := signals.StringSignals[key]; ok {
				return actual == expected
			}

			// int signal
			if actual, ok := signals.IntSignals[key]; ok {
				return actual == expected
			}

			// Signal doesn't exist - treat as false for bool, empty for string, 0 for int
			if expectedBool, ok := expected.(bool); ok {
				// If expecting false and signal doesn't exist, that's a match
				return !expectedBool
			}

		}
		return false
	}
}

// Allow external dynamic registration
func RegisterCondition(name string, fn ConditionFunc) {
	ConditionRegistry[name] = fn
}

// ===== Condition Evaluation Core =====
func evaluateCondition(raw interface{}, signals scanner.RepoSignals) bool {
	switch cond := raw.(type) {
	case map[string]interface{}:
		for key, val := range cond {
			if fn, ok := ConditionRegistry[key]; ok {
				return fn(val, signals)
			}
		}
		return false
	default:
		return false
	}
}

// ===== Rule Execution =====

func Evaluate(ruleSet []rules.Rule, signals scanner.RepoSignals) []Finding {
	var findings []Finding
	for _, r := range ruleSet {
		triggered := evaluateRule(r, signals)
		findings = append(findings, Finding{
			Rule:      r,
			Triggered: triggered,
		})
	}
	return findings
}

func evaluateRule(rule rules.Rule, signals scanner.RepoSignals) bool {
	for _, cond := range rule.Detect.NoneOf {
		if evaluateCondition(cond, signals) {
			return false
		}
	}
	for _, cond := range rule.Detect.AllOf {
		if !evaluateCondition(cond, signals) {
			return false
		}
	}
	if len(rule.Detect.AnyOf) > 0 {
		for _, cond := range rule.Detect.AnyOf {
			if evaluateCondition(cond, signals) {
				return true
			}
		}
		return false
	}
	return true
}
