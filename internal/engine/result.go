package engine

import "github.com/chuanjin/production-readiness/internal/rules"

// Finding represents the result of a single rule evaluation
type Finding struct {
	Rule      rules.Rule
	Triggered bool
}

// Summary aggregates findings and computes the score
type Summary struct {
	Total     int // Total number of rules evaluated
	Triggered int // Total number of rules triggered
	Passed    int // Total number of rules passed
	High      int // Number of high severity issues
	Medium    int // Number of medium severity issues
	Low       int // Number of low severity issues
	Score     int // Overall readiness score (0-100)
}

// Summarize calculates counts and a simple readiness score
func Summarize(findings []Finding) Summary {
	var s Summary

	s.Total = len(findings)

	for i := range findings {
		f := &findings[i]

		if !f.Triggered {
			s.Passed++
			continue
		}

		s.Triggered++

		switch f.Rule.Severity {
		case rules.High:
			s.High++
		case rules.Medium:
			s.Medium++
		case rules.Low:
			s.Low++
		}
	}

	// Simple scoring formula: 100 - (high*20 + medium*10 + low*5)
	// clamp to 0-100
	score := 100 - (s.High*20 + s.Medium*10 + s.Low*5)
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	s.Score = score

	return s
}

// GetSeverityCounts returns a breakdown of issues by severity
func (s Summary) GetSeverityCounts() map[string]int {
	return map[string]int{
		"high":   s.High,
		"medium": s.Medium,
		"low":    s.Low,
	}
}

// HasIssues returns true if any rules were triggered
func (s Summary) HasIssues() bool {
	return s.Triggered > 0
}

// IsProductionReady returns true if score is above threshold (default 80)
func (s Summary) IsProductionReady(threshold int) bool {
	if threshold == 0 {
		threshold = 80 // Default threshold
	}
	return s.Score >= threshold
}

// Grade returns a letter grade based on score
func (s Summary) Grade() string {
	switch {
	case s.Score >= 90:
		return "A"
	case s.Score >= 80:
		return "B"
	case s.Score >= 70:
		return "C"
	case s.Score >= 60:
		return "D"
	default:
		return "F"
	}
}
