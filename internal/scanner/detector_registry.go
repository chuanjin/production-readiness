// Package scanner
package scanner

// DetectorFunc is the signature for all detector functions
type DetectorFunc func(content string, relPath string, signals *RepoSignals)

// detectorRegistry holds all registered detectors
var detectorRegistry []DetectorFunc

// RegisterDetector adds a detector to the registry
// Call this from init() in detectors.go for each detector
func registerDetector(fn DetectorFunc) {
	detectorRegistry = append(detectorRegistry, fn)
}

// runAllDetectors executes all registered detectors
func runAllDetectors(content, relPath string, signals *RepoSignals) {
	for _, detector := range detectorRegistry {
		detector(content, relPath, signals)
	}
}
