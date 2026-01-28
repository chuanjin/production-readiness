// Package engine evaluates production-readiness rules against repository signals.
// It processes rule conditions (any_of, all_of, none_of) and determines which
// rules are triggered based on detected signals from the scanner.
package engine

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

// ConditionFunc evaluates a condition
type ConditionFunc func(value interface{}, signals *scanner.RepoSignals) bool

// ConditionRegistry holds all registered condition functions
var ConditionRegistry = map[string]ConditionFunc{}

func init() {
	// ===== Built-in evaluators =====

	ConditionRegistry["file_exists"] = func(value interface{}, signals *scanner.RepoSignals) bool {
		pattern := value.(string)

		// 1️⃣ basename exact match
		files := signals.GetFiles()
		for full := range files {
			if filepath.Base(full) == pattern {
				return true
			}
		}

		// 2️⃣ glob match using doublestar across full repo paths
		for full := range files {
			match, err := doublestar.Match(pattern, full)
			if err != nil {
				continue
			}
			if match {
				return true
			}
		}
		return false
	}

	ConditionRegistry["code_contains"] = func(value interface{}, signals *scanner.RepoSignals) bool {
		needle := value.(string)
		for _, content := range signals.GetFileContentMap() {
			if strings.Contains(content, needle) {
				return true
			}
		}
		return false
	}

	ConditionRegistry["signal_equals"] = func(value interface{}, signals *scanner.RepoSignals) bool {
		params := value.(map[string]interface{})
		for key, expected := range params {

			// bool signal
			if actual, ok := signals.GetBoolSignal(key); ok {
				return actual == expected
			}

			// string signal
			if actual, ok := signals.GetStringSignal(key); ok {
				return actual == expected
			}

			// int signal
			if actual, ok := signals.GetIntSignal(key); ok {
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

func RegisterCondition(name string, fn ConditionFunc) {
	ConditionRegistry[name] = fn
}

// ===== Condition Evaluation Core =====
func evaluateCondition(raw interface{}, signals *scanner.RepoSignals) bool {
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

func Evaluate(ruleSet []rules.Rule, signals *scanner.RepoSignals) []Finding {
	var findings []Finding
	// Use 'i' to avoid copying the 200-byte Rule struct into a local variable
	for i := range ruleSet {
		triggered := evaluateRule(&ruleSet[i], signals)
		findings = append(findings, Finding{
			Rule:      ruleSet[i], // This still copies into the new Finding
			Triggered: triggered,
		})
	}
	return findings
}

func evaluateRule(rule *rules.Rule, signals *scanner.RepoSignals) bool {
	// Evaluate all three condition groups independently
	noneOfPassed := evaluateNoneOf(rule.Detect.NoneOf, signals)
	allOfPassed := evaluateAllOf(rule.Detect.AllOf, signals)
	anyOfPassed := evaluateAnyOf(rule.Detect.AnyOf, signals)

	// Combine results with AND logic:
	// - none_of must pass (none of the conditions are true)
	// - all_of must pass (all conditions are true)
	// - any_of must pass (at least one condition is true, or no any_of exists)
	return noneOfPassed && allOfPassed && anyOfPassed
}

// evaluateNoneOf returns true if NONE of the conditions match
func evaluateNoneOf(conditions []map[string]interface{}, signals *scanner.RepoSignals) bool {
	// If no conditions, treat as passing (vacuous truth)
	if len(conditions) == 0 {
		return true
	}

	for _, cond := range conditions {
		if evaluateCondition(cond, signals) {
			return false // One matched, so none_of fails
		}
	}
	return true // None matched, so none_of passes
}

// evaluateAllOf returns true if ALL conditions match
func evaluateAllOf(conditions []map[string]interface{}, signals *scanner.RepoSignals) bool {
	// If no conditions, treat as passing (vacuous truth)
	if len(conditions) == 0 {
		return true
	}

	for _, cond := range conditions {
		if !evaluateCondition(cond, signals) {
			return false // One didn't match, so all_of fails
		}
	}
	return true // All matched, so all_of passes
}

// evaluateAnyOf returns true if at least ONE condition matches
// If no any_of conditions exist, returns true (vacuous truth)
func evaluateAnyOf(conditions []map[string]interface{}, signals *scanner.RepoSignals) bool {
	// If no conditions, treat as passing
	if len(conditions) == 0 {
		return true
	}

	for _, cond := range conditions {
		if evaluateCondition(cond, signals) {
			return true // One matched, so any_of passes
		}
	}
	return false // None matched, so any_of fails
}
