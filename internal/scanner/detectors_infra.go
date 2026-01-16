package scanner

import (
	"strings"
)

// detectSecretsProvider checks if code uses secrets management services
func detectSecretsProvider(content, relPath string, signals *RepoSignals) {
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

// detectInfrastructure checks if IaC (Infrastructure as Code) is present
func detectInfrastructure(content, relPath string, signals *RepoSignals) {
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
func detectRegions(content, relPath string, signals *RepoSignals) {
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
