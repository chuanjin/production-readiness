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
		{
			name:     "Invalid YAML content",
			content:  `invalid: yaml: content: [[[`,
			path:     "config.yaml",
			expected: false,
		},
		{
			name:     "Already detected signal",
			content:  `rate_limit: 100`,
			path:     "config.yaml",
			expected: true,
		},
		{
			name:     "Nested YAML rate limit",
			content:  `api:\n  gateway:\n    throttle:\n      enabled: true`,
			path:     "gateway.yml",
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

	// Test early return when signal already detected
	t.Run("Early return when already detected", func(t *testing.T) {
		signals := &RepoSignals{
			BoolSignals: map[string]bool{"api_gateway_rate_limit": true},
		}
		detectAPIGatewayRateLimit("rate_limit: 200", "config.yaml", signals)
		// Should still be true, function returns early
		if !signals.BoolSignals["api_gateway_rate_limit"] {
			t.Error("expected signal to remain true")
		}
	})
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
		{
			name:     "Weak indicators - multiple matches",
			content:  `99.9% availability four nines`,
			path:     "config.yaml",
			expected: true,
		},
		{
			name:     "Weak indicators - single match",
			content:  `99.9% uptime`,
			path:     "config.yaml",
			expected: false,
		},
		{
			name:     "ServiceLevelObjective kind",
			content:  "kind: ServiceLevelObjective\napiVersion: openslo/v1",
			path:     "slo.yaml",
			expected: true,
		},
		{
			name:     "Invalid YAML content",
			content:  `invalid: yaml: [[[`,
			path:     "config.yaml",
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

	// Test early return when signal already detected
	t.Run("Early return when already detected", func(t *testing.T) {
		signals := &RepoSignals{
			BoolSignals: map[string]bool{"slo_config_detected": true},
		}
		detectSLOConfig("service level objective: 99.9%", "readme.md", signals)
		// Should still be true, function returns early
		if !signals.BoolSignals["slo_config_detected"] {
			t.Error("expected signal to remain true")
		}
	})
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
		{
			name:     "Invalid YAML content",
			content:  `invalid: yaml: [[[`,
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

	// Test early return when signal already detected
	t.Run("Early return when already detected", func(t *testing.T) {
		signals := &RepoSignals{
			BoolSignals: map[string]bool{"error_budget_detected": true},
		}
		detectErrorBudget("error_budget: 0.1%", "config.yaml", signals)
		// Should still be true, function returns early
		if !signals.BoolSignals["error_budget_detected"] {
			t.Error("expected signal to remain true")
		}
	})
}

func TestCheckYAMLForErrorBudget(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected bool
	}{
		{
			name: "Error budget key exists",
			yaml: `
errorbudget: 0.1
`,
			expected: true,
		},
		{
			name: "Nested error budget key",
			yaml: `
monitoring:
  budget:
    burnrate: 14.4
`,
			expected: true,
		},
		{
			name: "List with error budget key",
			yaml: `
alerts:
  - name: my-alert
    burnrate: 10
`,
			expected: true,
		},
		{
			name: "Deeply nested array",
			yaml: `
services:
  - name: api
    monitoring:
      - type: slo
        errorbudget: 0.1
`,
			expected: true,
		},
		{
			name: "No error budget",
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
			result := checkYAMLForErrorBudget(doc)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDetectTimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		relPath  string
		expected bool
	}{
		{
			name: "Go HTTP client with timeout",
			content: `
				client := &http.Client{
					Timeout: 30 * time.Second,
				}
			`,
			relPath:  "main.go",
			expected: true,
		},
		{
			name: "Go context with timeout",
			content: `
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
			`,
			relPath:  "handler.go",
			expected: true,
		},
		{
			name: "Python requests with timeout",
			content: `
				response = requests.get('https://api.example.com', timeout=10)
			`,
			relPath:  "client.py",
			expected: true,
		},
		{
			name: "Node.js axios with timeout",
			content: `
				const client = axios.create({
					timeout: 5000,
					baseURL: 'https://api.example.com'
				});
			`,
			relPath:  "client.js",
			expected: true,
		},
		{
			name: "Database connection with timeout",
			content: `
				db.SetConnMaxLifetime(time.Minute * 3)
				db.SetMaxIdleConns(10)
			`,
			relPath:  "database.go",
			expected: true,
		},
		{
			name: "YAML config with timeout",
			content: `
				server:
				  port: 8080
				  timeout: 30s
				  read_timeout: 10s
			`,
			relPath:  "config.yaml",
			expected: true,
		},
		{
			name: "gRPC with timeout",
			content: `
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
				defer cancel()
				response, err := client.GetUser(ctx, request)
			`,
			relPath:  "grpc_client.go",
			expected: true,
		},
		{
			name: "No timeout configuration",
			content: `
				client := &http.Client{}
				resp, err := client.Get("https://api.example.com")
			`,
			relPath:  "main.go",
			expected: false,
		},
		{
			name: "Generic timeout keyword",
			content: `
				// Configure timeout for all operations
				const operationTimeout = 30
			`,
			relPath:  "config.go",
			expected: true,
		},
		{
			name: "Java HTTP client timeout",
			content: `
				HttpClient client = HttpClient.newBuilder()
					.connectTimeout(Duration.ofSeconds(10))
					.build();
			`,
			relPath:  "Client.java",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals:   make(map[string]bool),
				StringSignals: make(map[string]string),
			}
			detectTimeoutConfiguration(tt.content, tt.relPath, signals)

			if signals.BoolSignals["timeout_configured"] != tt.expected {
				t.Errorf("Expected timeout_configured=%v, got %v", tt.expected, signals.BoolSignals["timeout_configured"])
			}
		})
	}
}

func TestDetectRetry(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Go retry-go",
			content:  `retry.Do(func() error { return nil })`,
			expected: true,
		},
		{
			name:     "Python tenacity",
			content:  `@retry(stop=stop_after_attempt(3))`,
			expected: true,
		},
		{
			name:     "Generic retry limit",
			content:  `max_retries: 5`,
			expected: true,
		},
		{
			name:     "No retry",
			content:  `fmt.Println("hello")`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectRetry(tt.content, "test.go", signals)
			if signals.BoolSignals["retry_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["retry_detected"])
			}
		})
	}
}

func TestDetectCircuitBreaker(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Go gobreaker",
			content:  `cb := gobreaker.NewCircuitBreaker(st)`,
			expected: true,
		},
		{
			name:     "Java Resilience4j",
			content:  `CircuitBreaker circuitBreaker = CircuitBreaker.ofDefaults("backendName");`,
			expected: true,
		},
		{
			name:     "Generic circuit breaker",
			content:  `enable_circuit_breaker: true`,
			expected: true,
		},
		{
			name:     "No circuit breaker",
			content:  `func main() {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectCircuitBreaker(tt.content, "test.go", signals)
			if signals.BoolSignals["circuit_breaker_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["circuit_breaker_detected"])
			}
		})
	}
}
