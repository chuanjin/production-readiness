package engine

import (
	"reflect"
	"testing"

	"github.com/chuanjin/production-readiness/internal/rules"
)

func TestSummarize(t *testing.T) {
	tests := []struct {
		name     string
		findings []Finding
		expected Summary
	}{
		{
			name: "All passed",
			findings: []Finding{
				{Triggered: false, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: false, Rule: rules.Rule{Severity: rules.Medium}},
			},
			expected: Summary{
				Total:     2,
				Passed:    2,
				Triggered: 0,
				High:      0,
				Medium:    0,
				Low:       0,
				Score:     100,
			},
		},
		{
			name: "Mixed results",
			findings: []Finding{
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.Medium}},
				{Triggered: false, Rule: rules.Rule{Severity: rules.Low}},
			},
			expected: Summary{
				Total:     3,
				Passed:    1,
				Triggered: 2,
				High:      1,
				Medium:    1,
				Low:       0,
				Score:     70, // Score reduced by high (20) and medium (10) findings
			},
		},
		{
			name: "All triggered",
			findings: []Finding{
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
				{Triggered: true, Rule: rules.Rule{Severity: rules.High}},
			},
			expected: Summary{
				Total:     6,
				Passed:    0,
				Triggered: 6,
				High:      6,
				Medium:    0,
				Low:       0,
				Score:     0, // Score reduced to 0 by multiple high findings
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Summarize(tt.findings)
			if got != tt.expected {
				t.Errorf("Summarize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSummaryMethods(t *testing.T) {
	s := Summary{
		High:      2,
		Medium:    3,
		Low:       4,
		Score:     75,
		Triggered: 9,
	}

	t.Run("GetSeverityCounts", func(t *testing.T) {
		expected := map[string]int{
			"high":   2,
			"medium": 3,
			"low":    4,
		}
		if got := s.GetSeverityCounts(); !reflect.DeepEqual(got, expected) {
			t.Errorf("GetSeverityCounts() = %v, want %v", got, expected)
		}
	})

	t.Run("HasIssues", func(t *testing.T) {
		if !s.HasIssues() {
			t.Error("HasIssues() should be true")
		}
		empty := Summary{}
		if empty.HasIssues() {
			t.Error("HasIssues() should be false for empty summary")
		}
	})

	t.Run("IsProductionReady", func(t *testing.T) {
		if s.IsProductionReady(80) {
			t.Error("IsProductionReady(80) should be false for score 75")
		}
		if !s.IsProductionReady(70) {
			t.Error("IsProductionReady(70) should be true for score 75")
		}
		if s.IsProductionReady(0) { // Default 80
			t.Error("IsProductionReady(0) should be false for score 75 (default 80)")
		}
	})

	t.Run("Grade", func(t *testing.T) {
		tests := []struct {
			score int
			grade string
		}{
			{95, "A"},
			{85, "B"},
			{75, "C"},
			{65, "D"},
			{55, "F"},
		}

		for _, tt := range tests {
			s := Summary{Score: tt.score}
			if got := s.Grade(); got != tt.grade {
				t.Errorf("Grade() for score %d = %v, want %v", tt.score, got, tt.grade)
			}
		}
	})
}
