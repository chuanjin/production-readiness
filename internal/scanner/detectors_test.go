package scanner

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDetectSecretsProvider(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "AWS Secrets Manager detected",
			content:  `import "aws/secretsmanager"\nGetSecretValue()`,
			expected: true,
		},
		{
			name:     "HashiCorp Vault detected",
			content:  `import "hashicorp/vault"`,
			expected: true,
		},
		{
			name:     "No secrets provider",
			content:  `func main() { fmt.Println("hello") }`,
			expected: false,
		},
		{
			name:     "1Password detected",
			content:  `secret := "op://vault/item/field"`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectSecretsProvider(tt.content, "test.go", signals)

			if signals.BoolSignals["secrets_provider_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["secrets_provider_detected"])
			}
		})
	}
}

func TestDetectInfrastructure(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Terraform detected",
			content:  `provider "aws" { region = "us-east-1" }`,
			expected: true,
		},
		{
			name:     "Kubernetes YAML detected",
			content:  `apiVersion: v1\nkind: Deployment`,
			expected: true,
		},
		{
			name:     "No IaC detected",
			content:  `const x = 42;`,
			expected: false,
		},
		{
			name:     "CDK detected",
			content:  `import * as cdk from 'aws-cdk-lib';`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectInfrastructure(tt.content, "test.tf", signals)

			if signals.BoolSignals["infra_as_code_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["infra_as_code_detected"])
			}
		})
	}
}

func TestDetectRegions(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedCount int
	}{
		{
			name:          "Single AWS region",
			content:       `region = "us-east-1"`,
			expectedCount: 1,
		},
		{
			name:          "Multiple AWS regions",
			content:       `region = "us-east-1"\nbackup_region = "eu-west-1"`,
			expectedCount: 2,
		},
		{
			name:          "No regions",
			content:       `const x = 42;`,
			expectedCount: 0,
		},
		{
			name:          "GCP region",
			content:       `location = "us-central1"`,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				IntSignals: make(map[string]int),
			}
			detectRegions(tt.content, "test.tf", signals)

			if signals.IntSignals["region_count"] != tt.expectedCount {
				t.Errorf("expected %d regions, got %d", tt.expectedCount, signals.IntSignals["region_count"])
			}
		})
	}
}

func TestDetectVersionedArtifacts(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Semantic version detected",
			content:  `image: myapp:v1.2.3`,
			expected: true,
		},
		{
			name:     "SHA256 digest detected",
			content:  `image: myapp@sha256:abc123`,
			expected: true,
		},
		{
			name:     "Latest tag - not versioned",
			content:  `image: myapp:latest`,
			expected: false,
		},
		{
			name:     "Git tag detected",
			content:  `git tag v1.0.0`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectArtifactVersioning(tt.content, "deploy.yaml", signals)

			if signals.BoolSignals["versioned_artifacts"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["versioned_artifacts"])
			}
		})
	}
}

func TestDetectHealthEndpoints(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedEndpoint string
	}{
		{
			name:             "Health endpoint detected",
			content:          `app.get('/health', handler)`,
			expectedEndpoint: "/health",
		},
		{
			name:             "Ready endpoint detected",
			content:          `@route('/ready')\ndef ready():`,
			expectedEndpoint: "/ready",
		},
		{
			name:             "No endpoint",
			content:          `app.get('/api/users', handler)`,
			expectedEndpoint: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				StringSignals: make(map[string]string),
			}
			detectHealthEndpoints(tt.content, "server.js", signals)

			if signals.StringSignals["http_endpoint"] != tt.expectedEndpoint {
				t.Errorf("expected %q, got %q", tt.expectedEndpoint, signals.StringSignals["http_endpoint"])
			}
		})
	}
}

func TestDetectK8sProbes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Liveness probe detected",
			content: `
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        livenessProbe:
          httpGet:
            path: /health
`,
			expected: true,
		},
		{
			name: "Readiness probe detected",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    readinessProbe:
      httpGet:
        path: /ready
`,
			expected: true,
		},
		{
			name: "No probes",
			content: `
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectK8sProbes(tt.content, "deployment.yaml", signals)

			if signals.BoolSignals["k8s_probe_defined"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["k8s_probe_defined"])
			}
		})
	}
}

func TestDetectCorrelationId(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "X-Request-ID header",
			content:  `req.Header.Get("X-Request-ID")`,
			expected: true,
		},
		{
			name:     "Correlation ID variable",
			content:  `correlationId := generateID()`,
			expected: true,
		},
		{
			name:     "OpenTelemetry trace",
			content:  `import "go.opentelemetry.io/otel"`,
			expected: true,
		},
		{
			name:     "No correlation ID",
			content:  `log.Println("hello")`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectCorrelationID(tt.content, "handler.go", signals)

			if signals.BoolSignals["correlation_id_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["correlation_id_detected"])
			}
		})
	}
}

func TestDetectStructuredLogging(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Logrus detected",
			content:  `import "github.com/sirupsen/logrus"`,
			expected: true,
		},
		{
			name:     "Winston detected",
			content:  `const winston = require('winston');`,
			expected: true,
		},
		{
			name:     "Plain console.log",
			content:  `console.log("hello");`,
			expected: false,
		},
		{
			name:     "Zap logger",
			content:  `logger, _ := zap.NewProduction()`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectStructuredLogging(tt.content, "logger.go", signals)

			if signals.BoolSignals["structured_logging_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["structured_logging_detected"])
			}
		})
	}
}

func TestDetectMigrationTool(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Golang-migrate detected",
			content:  `import "github.com/golang-migrate/migrate"`,
			expected: true,
		},
		{
			name:     "Alembic detected",
			content:  `from alembic import context`,
			expected: true,
		},
		{
			name:     "Flyway detected",
			content:  `Flyway flyway = Flyway.configure()`,
			expected: true,
		},
		{
			name:     "No migration tool",
			content:  `SELECT * FROM users;`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectMigrationTool(tt.content, "migration.go", signals)

			if signals.BoolSignals["migration_tool_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["migration_tool_detected"])
			}
		})
	}
}

func TestDetectBackwardCompatibleMigration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Backward compatible mentioned",
			content:  `# This migration is backward compatible`,
			expected: true,
		},
		{
			name:     "Zero-downtime deployment",
			content:  `Supports zero-downtime deployment`,
			expected: true,
		},
		{
			name:     "Nullable column",
			content:  `ALTER TABLE users ADD COLUMN email VARCHAR(255) NULL DEFAULT '';`,
			expected: true,
		},
		{
			name:     "No compatibility info",
			content:  `ALTER TABLE users DROP COLUMN email;`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectBackwardCompatibleMigration(tt.content, "migration.sql", signals)

			if signals.BoolSignals["backward_compatible_migration_hint"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["backward_compatible_migration_hint"])
			}
		})
	}
}

func TestDetectMigrationValidation(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Dry-run flag",
			content:  `migrate up --dry-run`,
			expected: true,
		},
		{
			name:     "Validation test",
			content:  `def test_migration_validation():`,
			expected: true,
		},
		{
			name:     "Rollback test",
			content:  `it('should rollback successfully', () => {})`,
			expected: true,
		},
		{
			name:     "No validation",
			content:  `migrate up`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectMigrationValidation(tt.content, "test.sh", signals)

			if signals.BoolSignals["migration_validation_step"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["migration_validation_step"])
			}
		})
	}
}

func TestDetectResourceLimits(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Pod with limits detected",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    resources:
      limits:
        memory: "128Mi"
`,
			expected: true,
		},
		{
			name: "Deployment with limits detected",
			content: `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        resources:
          limits:
            cpu: "500m"
`,
			expected: true,
		},
		{
			name: "No limits defined",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    resources:
      requests:
        cpu: "100m"
`,
			expected: false,
		},
		{
			name: "No resources defined",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectResourceLimits(tt.content, "deploy.yaml", signals)

			if signals.BoolSignals["k8s_resource_limits_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["k8s_resource_limits_detected"])
			}
		})
	}
}

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
