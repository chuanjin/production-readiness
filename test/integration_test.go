package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

func TestEndToEndWorkflow(t *testing.T) {
	// Create a test repository
	tempDir := t.TempDir()

	// Scenario: Project with hardcoded secrets, no health checks, manual deployment
	testFiles := map[string]string{
		".env": `DATABASE_URL=postgres://user:pass@localhost/db
API_KEY=secret123`,
		"server.js": `const express = require('express');
const app = express();

app.get('/api/users', (req, res) => {
    res.json({ users: [] });
});

app.listen(3000);`,
		"README.md": `# Deployment Instructions

Step 1: SSH into the production server
Step 2: Run npm install
Step 3: Copy .env file manually
Step 4: Run pm2 restart app`,
		"docker-compose.yml": `version: '3'
services:
  app:
    image: myapp:latest
    ports:
      - "3000:3000"`,
	}

	for filename, content := range testFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Step 1: Scan the repository
	signals, err := scanner.ScanRepo(tempDir)
	if err != nil {
		t.Fatalf("ScanRepo failed: %v", err)
	}

	// Verify scanning worked
	if len(signals.Files) == 0 {
		t.Fatal("No files were scanned")
	}

	// Step 2: Load rules (create minimal test rules)
	testRules := []rules.Rule{
		{
			ID:       "hardcoded-secrets",
			Severity: "high",
			Category: "security",
			Title:    "Hardcoded secrets detected",
			Detect: rules.Detect{
				AnyOf: []map[string]interface{}{
					{"file_exists": ".env"},
					{"code_contains": "API_KEY"},
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
			ID:       "no-health-check",
			Severity: "medium",
			Category: "reliability",
			Title:    "No health check endpoint",
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
		{
			ID:       "manual-deployment",
			Severity: "medium",
			Category: "deployment",
			Title:    "Manual deployment detected",
			Detect: rules.Detect{
				AllOf: []map[string]interface{}{
					{
						"signal_equals": map[string]interface{}{
							"manual_steps_documented": true,
						},
					},
				},
			},
		},
		{
			ID:       "mutable-image-tag",
			Severity: "low",
			Category: "deployment",
			Title:    "Using mutable image tag",
			Detect: rules.Detect{
				NoneOf: []map[string]interface{}{
					{
						"signal_equals": map[string]interface{}{
							"versioned_artifacts": true,
						},
					},
				},
			},
		},
	}

	// Step 3: Evaluate rules
	findings := engine.Evaluate(testRules, signals)

	if len(findings) != 4 {
		t.Fatalf("Expected 4 findings, got %d", len(findings))
	}

	// Verify each expected finding
	expectedTriggers := map[string]bool{
		"hardcoded-secrets": true, // Should trigger (.env exists, no secrets provider)
		"no-health-check":   true, // Should trigger (no /health endpoint)
		"manual-deployment": true, // Should trigger (manual steps in README)
		"mutable-image-tag": true, // Should trigger (using :latest tag)
	}

	for _, finding := range findings {
		expected, exists := expectedTriggers[finding.Rule.ID]
		if !exists {
			t.Errorf("Unexpected rule ID: %s", finding.Rule.ID)
			continue
		}

		if finding.Triggered != expected {
			t.Errorf("Rule %s: expected triggered=%v, got %v",
				finding.Rule.ID, expected, finding.Triggered)
		}
	}

	// Count triggered findings
	triggeredCount := 0
	for _, finding := range findings {
		if finding.Triggered {
			triggeredCount++
		}
	}

	expectedTriggered := 4
	if triggeredCount != expectedTriggered {
		t.Errorf("Expected %d triggered findings, got %d", expectedTriggered, triggeredCount)
	}
}

func TestProductionReadyProject(t *testing.T) {
	// Create a well-configured project
	tempDir := t.TempDir()

	testFiles := map[string]string{
		"main.go": `package main

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func main() {
	log := logrus.WithFields(logrus.Fields{
		"request_id": "123",
	})
	log.Info("Starting server")
}`,
		"deployment.yaml": `apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  strategy:
    type: RollingUpdate
  template:
    spec:
      containers:
      - name: app
        image: myapp:v1.2.3
        livenessProbe:
          httpGet:
            path: /health
        readinessProbe:
          httpGet:
            path: /ready`,
		"secrets.yaml": `apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: app-secrets`,
		"terraform/main.tf": `provider "aws" {
  region = "us-east-1"
}

provider "aws" {
  alias  = "backup"
  region = "eu-west-1"
}`,
		"ingress.yaml": `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/limit-rps: "100"`,
	}

	for filename, content := range testFiles {
		path := filepath.Join(tempDir, filename)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	// Scan
	signals, err := scanner.ScanRepo(tempDir)
	if err != nil {
		t.Fatalf("ScanRepo failed: %v", err)
	}

	// Verify production-ready signals
	tests := []struct {
		name    string
		check   func() bool
		message string
	}{
		{
			name:    "Has structured logging",
			check:   func() bool { return signals.BoolSignals["structured_logging_detected"] },
			message: "Should detect logrus",
		},
		{
			name:    "Has correlation ID",
			check:   func() bool { return signals.BoolSignals["correlation_id_detected"] },
			message: "Should detect OpenTelemetry",
		},
		{
			name:    "Has health probes",
			check:   func() bool { return signals.BoolSignals["k8s_probe_defined"] },
			message: "Should detect Kubernetes probes",
		},
		{
			name:    "Has versioned artifacts",
			check:   func() bool { return signals.BoolSignals["versioned_artifacts"] },
			message: "Should detect v1.2.3 tag",
		},
		{
			name:    "Has secrets provider",
			check:   func() bool { return signals.BoolSignals["secrets_provider_detected"] },
			message: "Should detect external-secrets",
		},
		{
			name:    "Has multiple regions",
			check:   func() bool { return signals.IntSignals["region_count"] >= 2 },
			message: "Should detect 2 AWS regions",
		},
		{
			name:    "Has rate limiting",
			check:   func() bool { return signals.BoolSignals["ingress_rate_limit"] },
			message: "Should detect nginx rate limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.message)
			}
		})
	}
}
