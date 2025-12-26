package engine

import (
  "strings"

	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

type Finding struct {
	Rule      rules.Rule
	Triggered bool
}

type Summary struct {
	High     int
	Medium   int
	Low      int
	Positive int
	Score    int
}

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
	for _, cond := range rule.Detect.AnyOf {
		if matchCondition(cond, signals) {
			return true
		}
	}
	return false
}

func matchCondition(cond map[string]interface{}, signals scanner.RepoSignals) bool {
	for k, v := range cond {
		switch k {
		case "file_exists":
			name := v.(string)
			return signals.Files[name]
		case "code_contains":
			needle := v.(string)
			for _, c := range signals.FileContent {
				if strings.Contains(c, needle) {
					return true
				}
			}
		}
	}
	return false
}

func Summarize(findings []Finding) Summary {
	var sum Summary
	for _, f := range findings {
		if !f.Triggered {
			continue
		}
		switch f.Rule.Severity {
		case rules.High:
			sum.High++
		case rules.Medium:
			sum.Medium++
		case rules.Low:
			sum.Low++
		case rules.Positive:
			sum.Positive++
		}
	}
	score := 100 - (sum.High*20 + sum.Medium*10 + sum.Low*5)
	if score < 0 {
		score = 0
	}
	sum.Score = score
	return sum
}
