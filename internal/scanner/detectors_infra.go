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
	// Initialize if nil (defensive programming, though fs.go handles it)
	if signals.DetectedRegions == nil {
		signals.DetectedRegions = make(map[string]bool)
	}

	contentLower := strings.ToLower(content)

	// AWS regions
	awsRegions := []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"af-south-1", "ap-east-1", "ap-south-1", "ap-northeast-3", "ap-northeast-2",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
		"ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1",
		"eu-west-3", "eu-north-1", "me-south-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1",
	}

	// GCP regions
	gcpRegions := []string{
		"us-central1", "us-east1", "us-east4", "us-west1", "us-west2", "us-west3", "us-west4",
		"southamerica-east1", "northamerica-northeast1",
		"europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "europe-north1",
		"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3",
		"asia-southeast1", "asia-southeast2", "australia-southeast1",
	}

	// Azure regions
	azureRegions := []string{
		"eastus", "eastus2", "southcentralus", "westus2", "westus3", "australiaeast",
		"southeastasia", "northeurope", "westeurope", "uksouth", "ukwest", "francecentral",
		"germanywestcentral", "norwayeast", "switzerlandnorth", "japaneast", "japanwest",
		"centralindia", "southindia", "westindia", "canadacentral", "koreacentral",
	}

	allRegions := append(append(awsRegions, gcpRegions...), azureRegions...)

	for _, region := range allRegions {
		if strings.Contains(contentLower, region) {
			signals.DetectedRegions[region] = true
		}
	}

	// Update the global count based on the accumulated unique regions
	signals.IntSignals["region_count"] = len(signals.DetectedRegions)
}
