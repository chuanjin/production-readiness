package engine

import "github.com/chuanjin/production-readiness/internal/rules"

// Finding represents the result of a single rule evaluation
type Finding struct {
	Rule      rules.Rule
	Triggered bool
}

// Summary aggregates findings and computes the score
type Summary struct {
	High     int
	Medium   int
	Low      int
	Positive int
	Score    int
}

// Summarize calculates counts and a simple readiness score
func Summarize(findings []Finding) Summary {
	var s Summary

	for _, f := range findings {
		if !f.Triggered {
			continue
		}

		switch f.Rule.Severity {
		case rules.High:
			s.High++
		case rules.Medium:
			s.Medium++
		case rules.Low:
			s.Low++
		case rules.Positive:
			s.Positive++
		}
	}

	// Simple scoring formula: 100 - (high*20 + medium*10 + low*5) + positive*5
	// clamp to 0-100
	score := 100 - (s.High*20 + s.Medium*10 + s.Low*5) + (s.Positive * 5)
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	s.Score = score

	return s
}
