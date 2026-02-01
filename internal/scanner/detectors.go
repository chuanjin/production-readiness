package scanner

func init() {
	// Auto-register all detectors
	// Add new detectors here - just one line per detector!
	registerDetector(detectSecretsProvider)
	registerDetector(detectInfrastructure)
	registerDetector(detectRegions)
	registerDetector(detectNonRootUser)

	registerDetector(detectK8sDeploymentStrategy)
	registerDetector(detectK8sProbes)
	registerDetector(detectIngressRateLimit)
	registerDetector(detectResourceLimits)

	registerDetector(detectHealthEndpoints)
	registerDetector(detectCorrelationID)
	registerDetector(detectStructuredLogging)
	registerDetector(detectArtifactVersioning)

	registerDetector(detectAPIGatewayRateLimit)
	registerDetector(detectSLOConfig)
	registerDetector(detectErrorBudget)
	registerDetector(detectTimeoutConfiguration)
	registerDetector(detectRetry)
	registerDetector(detectCircuitBreaker)

	registerDetector(detectManualSteps)
	registerDetector(detectMigrationTool)
	registerDetector(detectBackwardCompatibleMigration)
	registerDetector(detectMigrationValidation)
	registerDetector(detectGracefulShutdown)
}

const (
	ExtYAML = ".yaml"
	ExtYML  = ".yml"
)
