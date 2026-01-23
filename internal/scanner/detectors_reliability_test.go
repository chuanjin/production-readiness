package scanner

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCheckYAMLForRateLimit(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected bool
	}{
		{
			name: "Rate limit key exists",
			yaml: `
rate_limit: 100
`,
			expected: true,
		},
		{
			name: "Nested rate limit key",
			yaml: `
services:
  api:
    throttle:
      burst: 10
`,
			expected: true,
		},
		{
			name: "List with rate limit key",
			yaml: `
plugins:
  - name: my-plugin
    rate-limit: 100
`,
			expected: true,
		},
		{
			name: "No rate limit",
			yaml: `
name: my-service
port: 8080
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc interface{}
			err := yaml.Unmarshal([]byte(tt.yaml), &doc)
			if err != nil {
				t.Fatalf("failed to unmarshal yaml: %v", err)
			}
			result := checkYAMLForRateLimit(doc)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDetectAPIGatewayRateLimit(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		path     string
		expected bool
	}{
		{
			name:     "AWS API Gateway throttle",
			content:  `"throttleSettings": { "burstLimit": 100 }`,
			path:     "template.json",
			expected: true,
		},
		{
			name:     "Nginx limit_req",
			content:  `limit_req zone=mylimit burst=20 nodelay;`,
			path:     "nginx.conf",
			expected: true,
		},
		{
			name:     "Generic rate limit",
			content:  `rate_limit: 100`,
			path:     "config.yaml",
			expected: true,
		},
		{
			name:     "No rate limit",
			content:  `server_name: example.com`,
			path:     "nginx.conf",
			expected: false,
		},
		{
			name:     "YAML structure rate limit",
			content:  `services:\n  api:\n    ratelimit:\n      max: 100`,
			path:     "config.yaml",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectAPIGatewayRateLimit(tt.content, tt.path, signals)

			if signals.BoolSignals["api_gateway_rate_limit"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["api_gateway_rate_limit"])
			}
		})
	}
}

func TestCheckYAMLForSLO(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected bool
	}{
		{
			name: "SLO key exists",
			yaml: `
slo: 99.9
`,
			expected: true,
		},
		{
			name: "Nested SLO key",
			yaml: `
monitoring:
  objective: 99.9
`,
			expected: true,
		},
		{
			name: "List with SLO key",
			yaml: `
rules:
  - name: my-rule
    slo: 99.9
`,
			expected: true,
		},
		{
			name: "No SLO",
			yaml: `
name: my-service
port: 8080
`,
			expected: false,
		},
		{
			name: "Case insensitive",
			yaml: `
SLO: 99.9
`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc interface{}
			err := yaml.Unmarshal([]byte(tt.yaml), &doc)
			if err != nil {
				t.Fatalf("failed to unmarshal yaml: %v", err)
			}
			result := checkYAMLForSLO(doc)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDetectSLOConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		path     string
		expected bool
	}{
		{
			name:     "OpenSLO kind",
			content:  "kind: SLO\napiVersion: openslo/v1",
			path:     "slo.yaml",
			expected: true,
		},
		{
			name:     "Slo keywords",
			content:  `service level objective: 99.9%`,
			path:     "readme.md",
			expected: true,
		},
		{
			name: "YAML structure SLO",
			content: `monitoring:
  slo:
    availability: 99.9`,
			path:     "monitor.yaml",
			expected: true,
		},
		{
			name:     "No SLO",
			content:  "apiVersion: v1\nkind: Pod",
			path:     "pod.yaml",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectSLOConfig(tt.content, tt.path, signals)

			if signals.BoolSignals["slo_config_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["slo_config_detected"])
			}
		})
	}
}

func TestDetectErrorBudget(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		path     string
		expected bool
	}{
		{
			name:     "Error budget keyword",
			content:  `error_budget: 0.1%`,
			path:     "config.yaml",
			expected: true,
		},
		{
			name:     "Burn rate",
			content:  `alert_if: burn_rate > 10`,
			path:     "alerts.yaml",
			expected: true,
		},
		{
			name: "YAML structure error budget",
			content: `monitoring:
  budget:
    burnrate: 14.4`,
			path:     "monitor.yaml",
			expected: true,
		},
		{
			name:     "No error budget",
			content:  `log_level: info`,
			path:     "config.yaml",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectErrorBudget(tt.content, tt.path, signals)

			if signals.BoolSignals["error_budget_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["error_budget_detected"])
			}
		})
	}
}
