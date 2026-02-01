package scanner

import (
	"strings"

	"github.com/chuanjin/production-readiness/internal/patterns"
)

// detectSecretsProvider checks if code uses secrets management services
func detectSecretsProvider(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("secrets_provider_detected") {
		return
	}

	secretsProviderPatterns := patterns.SecretsProviderPatterns

	contentLower := strings.ToLower(content)
	for _, pattern := range secretsProviderPatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			signals.SetBool("secrets_provider_detected", true)
			return
		}
	}
}

// detectInfrastructure checks if IaC (Infrastructure as Code) is present
func detectInfrastructure(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("infra_as_code_detected") {
		return
	}

	contentLower := strings.ToLower(content)

	infraPatterns := patterns.InfraPatterns

	for _, pattern := range infraPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.SetBool("infra_as_code_detected", true)
			return
		}
	}
}

// detectRegions counts the number of unique cloud regions configured
func detectRegions(content, relPath string, signals *RepoSignals) {
	contentLower := strings.ToLower(content)

	// AWS regions
	awsRegions := patterns.AWSRegions

	// GCP regions
	gcpRegions := patterns.GCPRegions

	// Azure regions
	azureRegions := patterns.AzureRegions

	allRegions := append(append(awsRegions, gcpRegions...), azureRegions...)

	for _, region := range allRegions {
		if strings.Contains(contentLower, region) {
			signals.SetRegion(region)
		}
	}

	// Update the global count based on the accumulated unique regions
	signals.SetInt("region_count", signals.GetRegionCount())
}

// detectNonRootUser checks for non-root user configuration in Dockerfiles
func detectNonRootUser(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("non_root_user_detected") {
		return
	}

	// Only scan Dockerfiles
	fileName := strings.ToLower(relPath)
	if !strings.Contains(fileName, "dockerfile") {
		return
	}

	contentLower := strings.ToLower(content)
	nonRootUserPatterns := patterns.NonRootUserPatterns

	for _, pattern := range nonRootUserPatterns {
		if strings.Contains(contentLower, pattern) {
			// Basic check: Ensure it's not USER root
			lines := strings.Split(contentLower, "\n")
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLine, "user ") || strings.HasPrefix(trimmedLine, "user\t") {
					parts := strings.Fields(trimmedLine)
					if len(parts) >= 2 && parts[1] != "root" && parts[1] != "0" {
						signals.SetBool("non_root_user_detected", true)
						return
					}
				}
			}
		}
	}
}
