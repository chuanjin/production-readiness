package scanner

// DetectorFunc is the signature for all detector functions
type DetectorFunc func(content string, relPath string, signals *RepoSignals)

// detectorRegistry holds all registered detectors
var detectorRegistry []DetectorFunc

func init() {
	// Register all detectors here
	// Add new detectors to this list and they'll run automatically
	detectorRegistry = []DetectorFunc{
		detectSecretsProvider,
		detectInfrastructure,
		detectRegions,
		detectManualSteps,
		detectK8sDeploymentStrategy,
		detectArtifactVersioning,
	}
}

// runAllDetectors executes all registered detectors
func runAllDetectors(content string, relPath string, signals *RepoSignals) {
	for _, detector := range detectorRegistry {
		detector(content, relPath, signals)
	}
}
