package patterns

import (
	"testing"
)

func TestPatterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{"APIGatewayRateLimitPatterns", APIGatewayRateLimitPatterns},
		{"SLOPatterns", SLOPatterns},
		{"ErrorBudgetPatterns", ErrorBudgetPatterns},
		{"SecretsProviderPatterns", SecretsProviderPatterns},
		{"InfraPatterns", InfraPatterns},
		{"AWSRegions", AWSRegions},
		{"GCPRegions", GCPRegions},
		{"AzureRegions", AzureRegions},
		{"ManualStepPatterns", ManualStepPatterns},
		{"MigrationToolPatterns", MigrationToolPatterns},
		{"BackwardCompatPatterns", BackwardCompatPatterns},
		{"MigrationValidationPatterns", MigrationValidationPatterns},
		{"MutableTags", MutableTags},
		{"VersioningPatterns", VersioningPatterns},
		{"HealthPatterns", HealthPatterns},
		{"ReadyPatterns", ReadyPatterns},
		{"CorrelationPatterns", CorrelationPatterns},
		{"StructuredLoggingPatterns", StructuredLoggingPatterns},
		{"StrongStructuredLoggingIndicators", StrongStructuredLoggingIndicators},
		{"NginxIngressRateLimitAnnotations", NginxIngressRateLimitAnnotations},
		{"RateLimitYAMLKeys", RateLimitYAMLKeys},
		{"SLOYAMLKeys", SLOYAMLKeys},
		{"ErrorBudgetYAMLKeys", ErrorBudgetYAMLKeys},
		{"DocFileKeywords", DocFileKeywords},
		{"DocFileExtensions", DocFileExtensions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.patterns) == 0 {
				t.Errorf("%s is empty", tt.name)
			}

			seen := make(map[string]bool)
			for _, p := range tt.patterns {
				if seen[p] {
					t.Errorf("%s contains duplicate pattern: %s", tt.name, p)
				}
				seen[p] = true
			}
		})
	}
}

func TestK8sValidKinds(t *testing.T) {
	if len(K8sValidKinds) == 0 {
		t.Error("K8sValidKinds is empty")
	}
	// Verify some expected keys exist
	expected := []string{"Pod", "Deployment", "StatefulSet"}
	for _, k := range expected {
		if !K8sValidKinds[k] {
			t.Errorf("K8sValidKinds missing expected kind: %s", k)
		}
	}
}
