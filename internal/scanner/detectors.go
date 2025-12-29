package scanner

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// detectSecretsProvider checks if code uses secrets management services
func detectSecretsProvider(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["secrets_provider_detected"] {
		return
	}

	secretsProviderPatterns := []string{
		// AWS Secrets Manager
		"aws-sdk", "aws/secretsmanager", "GetSecretValue", "secretsmanager",
		"AWS::SecretsManager",

		// HashiCorp Vault
		"hashicorp/vault", "vault.NewClient", "vault/api",

		// Google Secret Manager
		"cloud.google.com/go/secretmanager", "secretmanager.NewClient",
		"google-cloud/secret-manager",

		// Azure Key Vault
		"azure-keyvault", "azure/keyvault", "KeyVaultClient",

		// Doppler
		"doppler.com", "DopplerSDK", "@dopplerhq",

		// Infisical
		"infisical", "infisical-sdk",

		// 1Password
		"1password", "op://",

		// Generic secrets management
		"sealed-secrets", "external-secrets", "secrets-store-csi",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range secretsProviderPatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			signals.BoolSignals["secrets_provider_detected"] = true
			return
		}
	}
}

// detectK8sDeploymentStrategy checks Kubernetes deployment files for strategy
func detectK8sDeploymentStrategy(content string, relPath string, signals *RepoSignals) {
	if signals.StringSignals["k8s_deployment_strategy"] != "" {
		return
	}

	// Only check YAML files
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext != ".yaml" && ext != ".yml" {
		return
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return
	}

	kind, ok := doc["kind"].(string)
	if !ok || kind != "Deployment" {
		return
	}

	if spec, ok := doc["spec"].(map[string]interface{}); ok {
		if strategy, ok := spec["strategy"].(map[string]interface{}); ok {
			if strategyType, ok := strategy["type"].(string); ok {
				signals.StringSignals["k8s_deployment_strategy"] = strategyType
			}
		}
	}
}

// detectArtifactVersioning checks for versioned artifact patterns
func detectArtifactVersioning(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["versioned_artifacts"] {
		return
	}

	contentLower := strings.ToLower(content)

	// Check for mutable tags first (anti-pattern)
	mutableTags := []string{":latest", ":main", ":master", ":dev", ":develop"}
	for _, tag := range mutableTags {
		if strings.Contains(contentLower, tag) {
			// Found mutable tag - not versioned
			return
		}
	}

	// Look for versioning patterns
	versioningPatterns := []string{
		// Semantic versioning
		":v1", ":v2", "version:", "tag:",

		// Git tags
		"git tag", "github.ref", "git.tag",

		// Semantic versioning tools
		"semver", "semantic-release",

		// Docker image versioning
		"@sha256:", "sha-", ":build-", ":release-",

		// Container registries with versions
		"gcr.io", "ecr.aws", "quay.io", "ghcr.io",

		// Version variables
		"$version", "${version}", "{{version}}",
	}

	for _, pattern := range versioningPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["versioned_artifacts"] = true
			return
		}
	}
}

// detectInfrastructure checks if IaC (Infrastructure as Code) is present
func detectInfrastructure(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["infra_as_code_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	infraPatterns := []string{
		// Terraform
		"terraform", "provider \"", "resource \"", "module \"",

		// CloudFormation
		"aws::cloudformation", "awscloudformation", "resources:",

		// Pulumi
		"pulumi", "@pulumi/",

		// CDK
		"aws-cdk", "@aws-cdk/",

		// Kubernetes/Helm
		"apiversion:", "kind: deployment", "kind: service",

		// Ansible
		"ansible", "playbook",
	}

	for _, pattern := range infraPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["infra_as_code_detected"] = true
			return
		}
	}
}

// detectRegions counts the number of unique cloud regions configured
func detectRegions(content string, relPath string, signals *RepoSignals) {
	regions := make(map[string]bool)

	contentLower := strings.ToLower(content)

	// AWS regions
	awsRegions := []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-north-1",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ap-northeast-2",
		"sa-east-1", "ca-central-1", "ap-south-1",
	}

	// GCP regions
	gcpRegions := []string{
		"us-central1", "us-east1", "us-west1", "us-east4",
		"europe-west1", "europe-west2", "europe-west3", "europe-north1",
		"asia-east1", "asia-southeast1", "asia-northeast1",
	}

	// Azure regions
	azureRegions := []string{
		"eastus", "eastus2", "westus", "westus2", "centralus",
		"northeurope", "westeurope", "uksouth", "ukwest",
		"southeastasia", "eastasia", "japaneast", "japanwest",
	}

	allRegions := append(append(awsRegions, gcpRegions...), azureRegions...)

	for _, region := range allRegions {
		if strings.Contains(contentLower, region) {
			regions[region] = true
		}
	}

	// Update the count (only if we found more regions than before)
	if len(regions) > signals.IntSignals["region_count"] {
		signals.IntSignals["region_count"] = len(regions)
	}
}

// detectManualSteps checks if documentation contains manual deployment steps
func detectManualSteps(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["manual_steps_documented"] {
		return
	}

	// Only check documentation files
	fileName := strings.ToLower(relPath)
	isDocFile := strings.Contains(fileName, "readme") ||
		strings.Contains(fileName, "doc") ||
		strings.Contains(fileName, "deploy") ||
		strings.Contains(fileName, "setup") ||
		strings.Contains(fileName, "install") ||
		strings.HasSuffix(fileName, ".md") ||
		strings.HasSuffix(fileName, ".txt")

	if !isDocFile {
		return
	}

	contentLower := strings.ToLower(content)

	// Patterns indicating manual steps
	manualStepPatterns := []string{
		// Step-by-step instructions
		"step 1", "step 2", "1.", "2.", "3.",
		"first,", "then,", "next,", "finally,",

		// Manual actions
		"manually", "by hand", "login to", "navigate to",
		"click on", "open the", "go to the console",
		"ssh into", "copy the file", "run this command",

		// Console/UI instructions
		"in the console", "in the dashboard", "in the ui",
		"from the web interface", "using the portal",

		// Manual verification
		"verify that", "check that", "make sure",
		"confirm that", "ensure that",

		// Manual configuration
		"edit the file", "update the", "change the",
		"set the value", "configure manually",
	}

	// Count matches to avoid false positives
	matches := 0
	for _, pattern := range manualStepPatterns {
		if strings.Contains(contentLower, pattern) {
			matches++
			if matches >= 3 { // Need at least 3 indicators
				signals.BoolSignals["manual_steps_documented"] = true
				return
			}
		}
	}
}
