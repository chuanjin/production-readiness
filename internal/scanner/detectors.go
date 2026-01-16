package scanner

func init() {
	// Auto-register all detectors
	// Add new detectors here - just one line per detector!
	registerDetector(detectSecretsProvider)
	registerDetector(detectInfrastructure)
	registerDetector(detectRegions)

	registerDetector(detectK8sDeploymentStrategy)
	registerDetector(detectK8sProbes)
	registerDetector(detectIngressRateLimit)

	registerDetector(detectHealthEndpoints)
	registerDetector(detectCorrelationID)
	registerDetector(detectStructuredLogging)
	registerDetector(detectArtifactVersioning)

	registerDetector(detectAPIGatewayRateLimit)
	registerDetector(detectSLOConfig)
	registerDetector(detectErrorBudget)

	registerDetector(detectManualSteps)
	registerDetector(detectMigrationTool)
	registerDetector(detectBackwardCompatibleMigration)
	registerDetector(detectMigrationValidation)
}

const (
	ExtYAML = ".yaml"
	ExtYML  = ".yml"
)
