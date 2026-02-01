package scanner

import (
	"testing"
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
	// 1. Test basic detection and global accumulation across multiple calls
	t.Run("Global accumulation", func(t *testing.T) {
		signals := &RepoSignals{
			IntSignals:      make(map[string]int),
			DetectedRegions: make(map[string]bool),
		}

		// First file: us-east-1
		detectRegions(`region = "us-east-1"`, "file1.tf", signals)
		if signals.IntSignals["region_count"] != 1 {
			t.Errorf("expected 1 region, got %d", signals.IntSignals["region_count"])
		}

		// Second file: eu-west-1 (should increment count)
		detectRegions(`backup_region = "eu-west-1"`, "file2.tf", signals)
		if signals.IntSignals["region_count"] != 2 {
			t.Errorf("expected 2 regions, got %d", signals.IntSignals["region_count"])
		}

		// Third file: us-east-1 again (duplicate, should NOT increment)
		detectRegions(`another_ref = "us-east-1"`, "file3.tf", signals)
		if signals.IntSignals["region_count"] != 2 {
			t.Errorf("expected 2 regions, got %d", signals.IntSignals["region_count"])
		}
	})

	// 2. Test specific new regions
	tests := []struct {
		name          string
		content       string
		expectedCount int
	}{
		{
			name:          "New AWS region (af-south-1)", // Cape Town
			content:       `region = "af-south-1"`,
			expectedCount: 1,
		},
		{
			name:          "New GCP region (southamerica-east1)",
			content:       `location = "southamerica-east1"`,
			expectedCount: 1,
		},
		{
			name:          "New Azure region (switzerlandnorth)",
			content:       `location = "switzerlandnorth"`,
			expectedCount: 1,
		},
		{
			name:          "Multiple mixed regions",
			content:       `"us-east-1", "europe-west1", "japaneast"`,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				IntSignals:      make(map[string]int),
				DetectedRegions: make(map[string]bool),
			}
			detectRegions(tt.content, "test.tf", signals)

			if signals.IntSignals["region_count"] != tt.expectedCount {
				t.Errorf("expected %d regions, got %d", tt.expectedCount, signals.IntSignals["region_count"])
			}
		})
	}
}

func TestDetectNonRootUser(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		relPath  string
		expected bool
	}{
		{
			name: "Non-root user detected",
			content: `
FROM golang:1.21
RUN useradd -m myuser
USER myuser
ENTRYPOINT ["./app"]
`,
			relPath:  "Dockerfile",
			expected: true,
		},
		{
			name: "UID user detected",
			content: `
FROM alpine
USER 1000
`,
			relPath:  "dockerfile",
			expected: true,
		},
		{
			name: "Explicit root user (root)",
			content: `
FROM ubuntu
USER root
`,
			relPath:  "Dockerfile",
			expected: false,
		},
		{
			name: "Explicit root user (0)",
			content: `
FROM ubuntu
USER 0
`,
			relPath:  "Dockerfile",
			expected: false,
		},
		{
			name: "No USER instruction",
			content: `
FROM node:18
COPY . .
`,
			relPath:  "Dockerfile",
			expected: false,
		},
		{
			name: "Not a Dockerfile",
			content: `
USER myuser
`,
			relPath:  "readme.md",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectNonRootUser(tt.content, tt.relPath, signals)

			if signals.BoolSignals["non_root_user_detected"] != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.BoolSignals["non_root_user_detected"])
			}
		})
	}
}
