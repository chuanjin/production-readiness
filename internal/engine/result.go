package engine

import "github.com/chuanjin/production-readiness/internal/rules"

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
