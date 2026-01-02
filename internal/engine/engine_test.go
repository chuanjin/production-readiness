package engine

import (
	"testing"

	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

func TestSummarize_ScoreCalculation(t *testing.T) {
	// 1 high, 2 medium, 3 low triggered
	findings := []Finding{
		{Rule: rules.Rule{Severity: rules.High}, Triggered: true},
		{Rule: rules.Rule{Severity: rules.Medium}, Triggered: true},
		{Rule: rules.Rule{Severity: rules.Medium}, Triggered: true},
		{Rule: rules.Rule{Severity: rules.Low}, Triggered: true},
		{Rule: rules.Rule{Severity: rules.Low}, Triggered: true},
		{Rule: rules.Rule{Severity: rules.Low}, Triggered: true},
	}

	s := Summarize(findings)
	// formula: 100 - (high*20 + medium*10 + low*5)
	exp := 100 - (1*20 + 2*10 + 3*5)
	if s.Score != exp {
		t.Fatalf("expected score %d, got %d", exp, s.Score)
	}
}

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		name      string
		condition interface{}
		signals   scanner.RepoSignals
		expected  bool
	}{
		{
			name: "file_exists - exact match",
			condition: map[string]interface{}{
				"file_exists": ".env",
			},
			signals: scanner.RepoSignals{
				Files: map[string]bool{
					".env":    true,
					"main.go": true,
				},
			},
			expected: true,
		},
		{
			name: "file_exists - not found",
			condition: map[string]interface{}{
				"file_exists": ".env",
			},
			signals: scanner.RepoSignals{
				Files: map[string]bool{
					"main.go": true,
				},
			},
			expected: false,
		},
		{
			name: "file_exists - glob pattern",
			condition: map[string]interface{}{
				"file_exists": "**/*.yaml",
			},
			signals: scanner.RepoSignals{
				Files: map[string]bool{
					"configs/deployment.yaml": true,
					"main.go":                 true,
				},
			},
			expected: true,
		},
		{
			name: "code_contains - found",
			condition: map[string]interface{}{
				"code_contains": "process.env",
			},
			signals: scanner.RepoSignals{
				FileContent: map[string]string{
					"server.js": "const port = process.env.PORT",
				},
			},
			expected: true,
		},
		{
			name: "code_contains - not found",
			condition: map[string]interface{}{
				"code_contains": "process.env",
			},
			signals: scanner.RepoSignals{
				FileContent: map[string]string{
					"server.js": "const port = 3000",
				},
			},
			expected: false,
		},
		{
			name: "signal_equals - bool true",
			condition: map[string]interface{}{
				"signal_equals": map[string]interface{}{
					"secrets_provider_detected": true,
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"secrets_provider_detected": true,
				},
			},
			expected: true,
		},
		{
			name: "signal_equals - bool false",
			condition: map[string]interface{}{
				"signal_equals": map[string]interface{}{
					"secrets_provider_detected": false,
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"secrets_provider_detected": true,
				},
			},
			expected: false,
		},
		{
			name: "signal_equals - bool not set (treated as false)",
			condition: map[string]interface{}{
				"signal_equals": map[string]interface{}{
					"secrets_provider_detected": false,
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{},
			},
			expected: true,
		},
		{
			name: "signal_equals - string",
			condition: map[string]interface{}{
				"signal_equals": map[string]interface{}{
					"http_endpoint": "/health",
				},
			},
			signals: scanner.RepoSignals{
				StringSignals: map[string]string{
					"http_endpoint": "/health",
				},
			},
			expected: true,
		},
		{
			name: "signal_equals - int",
			condition: map[string]interface{}{
				"signal_equals": map[string]interface{}{
					"region_count": 1,
				},
			},
			signals: scanner.RepoSignals{
				IntSignals: map[string]int{
					"region_count": 1,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize signal maps if nil
			if tt.signals.Files == nil {
				tt.signals.Files = make(map[string]bool)
			}
			if tt.signals.FileContent == nil {
				tt.signals.FileContent = make(map[string]string)
			}
			if tt.signals.BoolSignals == nil {
				tt.signals.BoolSignals = make(map[string]bool)
			}
			if tt.signals.StringSignals == nil {
				tt.signals.StringSignals = make(map[string]string)
			}
			if tt.signals.IntSignals == nil {
				tt.signals.IntSignals = make(map[string]int)
			}

			result := evaluateCondition(tt.condition, tt.signals)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateRule(t *testing.T) {
	tests := []struct {
		name     string
		rule     rules.Rule
		signals  scanner.RepoSignals
		expected bool
	}{
		{
			name: "none_of - should trigger when condition is false",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					NoneOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"secrets_provider_detected": true,
							},
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"secrets_provider_detected": false,
				},
			},
			expected: true,
		},
		{
			name: "none_of - should not trigger when condition is true",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					NoneOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"secrets_provider_detected": true,
							},
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"secrets_provider_detected": true,
				},
			},
			expected: false,
		},
		{
			name: "all_of - should trigger when all conditions are true",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					AllOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"infra_as_code_detected": true,
							},
						},
						{
							"signal_equals": map[string]interface{}{
								"region_count": 1,
							},
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"infra_as_code_detected": true,
				},
				IntSignals: map[string]int{
					"region_count": 1,
				},
			},
			expected: true,
		},
		{
			name: "all_of - should not trigger when one condition is false",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					AllOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"infra_as_code_detected": true,
							},
						},
						{
							"signal_equals": map[string]interface{}{
								"region_count": 1,
							},
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"infra_as_code_detected": true,
				},
				IntSignals: map[string]int{
					"region_count": 2, // Different value
				},
			},
			expected: false,
		},
		{
			name: "any_of - should trigger when at least one condition is true",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					AnyOf: []map[string]interface{}{
						{
							"file_exists": ".env",
						},
						{
							"code_contains": "process.env",
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				Files: map[string]bool{
					".env": true,
				},
				FileContent: map[string]string{},
			},
			expected: true,
		},
		{
			name: "any_of - should not trigger when all conditions are false",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					AnyOf: []map[string]interface{}{
						{
							"file_exists": ".env",
						},
						{
							"code_contains": "process.env",
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				Files:       map[string]bool{},
				FileContent: map[string]string{},
			},
			expected: false,
		},
		{
			name: "complex rule - any_of + none_of",
			rule: rules.Rule{
				ID: "test-rule",
				Detect: rules.Detect{
					AnyOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"migration_tool_detected": true,
							},
						},
					},
					NoneOf: []map[string]interface{}{
						{
							"signal_equals": map[string]interface{}{
								"backward_compatible_migration_hint": true,
							},
						},
					},
				},
			},
			signals: scanner.RepoSignals{
				BoolSignals: map[string]bool{
					"migration_tool_detected":            true,
					"backward_compatible_migration_hint": false,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize signal maps if nil
			if tt.signals.Files == nil {
				tt.signals.Files = make(map[string]bool)
			}
			if tt.signals.FileContent == nil {
				tt.signals.FileContent = make(map[string]string)
			}
			if tt.signals.BoolSignals == nil {
				tt.signals.BoolSignals = make(map[string]bool)
			}
			if tt.signals.StringSignals == nil {
				tt.signals.StringSignals = make(map[string]string)
			}
			if tt.signals.IntSignals == nil {
				tt.signals.IntSignals = make(map[string]int)
			}

			result := evaluateRule(tt.rule, tt.signals)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	rules := []rules.Rule{
		{
			ID:       "secrets-hardcoded",
			Severity: "high",
			Detect: rules.Detect{
				AnyOf: []map[string]interface{}{
					{"file_exists": ".env"},
				},
				NoneOf: []map[string]interface{}{
					{
						"signal_equals": map[string]interface{}{
							"secrets_provider_detected": true,
						},
					},
				},
			},
		},
		{
			ID:       "health-check",
			Severity: "medium",
			Detect: rules.Detect{
				NoneOf: []map[string]interface{}{
					{
						"signal_equals": map[string]interface{}{
							"http_endpoint": "/health",
						},
					},
				},
			},
		},
	}

	signals := scanner.RepoSignals{
		Files: map[string]bool{
			".env": true,
		},
		FileContent: map[string]string{},
		BoolSignals: map[string]bool{
			"secrets_provider_detected": false,
		},
		StringSignals: map[string]string{
			"http_endpoint": "/health",
		},
		IntSignals: make(map[string]int),
	}

	findings := Evaluate(rules, signals)

	if len(findings) != 2 {
		t.Fatalf("Expected 2 findings, got %d", len(findings))
	}

	// First rule should trigger (has .env but no secrets provider)
	if !findings[0].Triggered {
		t.Error("secrets-hardcoded rule should have triggered")
	}
	if findings[0].Rule.ID != "secrets-hardcoded" {
		t.Errorf("Expected rule ID 'secrets-hardcoded', got %q", findings[0].Rule.ID)
	}

	// Second rule should NOT trigger (has /health endpoint)
	if findings[1].Triggered {
		t.Error("health-check rule should not have triggered")
	}
	if findings[1].Rule.ID != "health-check" {
		t.Errorf("Expected rule ID 'health-check', got %q", findings[1].Rule.ID)
	}
}
