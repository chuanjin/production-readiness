package scanner

import (
	"strings"

	"github.com/chuanjin/production-readiness/internal/patterns"
)

// detectArtifactVersioning checks for versioned artifact patterns
func detectArtifactVersioning(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["versioned_artifacts"] {
		return
	}

	contentLower := strings.ToLower(content)

	// Check for mutable tags first (anti-pattern)
	mutableTags := patterns.MutableTags
	for _, tag := range mutableTags {
		if strings.Contains(contentLower, tag) {
			// Found mutable tag - not versioned
			return
		}
	}

	// Look for versioning patterns
	versioningPatterns := patterns.VersioningPatterns

	for _, pattern := range versioningPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["versioned_artifacts"] = true
			return
		}
	}
}

// detectHealthEndpoints checks for health check HTTP endpoints
func detectHealthEndpoints(content, relPath string, signals *RepoSignals) {
	contentLower := strings.ToLower(content)

	// Detect /health endpoint
	if signals.StringSignals["http_endpoint"] == "" {
		healthPatterns := patterns.HealthPatterns

		for _, pattern := range healthPatterns {
			if strings.Contains(contentLower, pattern) {
				signals.StringSignals["http_endpoint"] = "/health"
				break
			}
		}
	}

	// Detect /ready or /readiness endpoint
	if signals.StringSignals["http_endpoint"] == "" {
		readyPatterns := patterns.ReadyPatterns

		for _, pattern := range readyPatterns {
			if strings.Contains(contentLower, pattern) {
				signals.StringSignals["http_endpoint"] = "/ready"
				break
			}
		}
	}
}

// detectCorrelationID checks for correlation/trace ID usage
func detectCorrelationID(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["correlation_id_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	correlationPatterns := patterns.CorrelationPatterns

	for _, pattern := range correlationPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["correlation_id_detected"] = true
			return
		}
	}
}

// detectStructuredLogging checks for structured logging libraries and patterns
func detectStructuredLogging(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["structured_logging_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	structuredLoggingPatterns := patterns.StructuredLoggingPatterns

	matchCount := 0
	for _, pattern := range structuredLoggingPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Need at least 2 matches to be confident it's structured logging
			// (to avoid false positives from just having "log.info")
			if matchCount >= 2 {
				signals.BoolSignals["structured_logging_detected"] = true
				return
			}
		}
	}

	// Single strong indicator is enough
	strongIndicators := patterns.StrongStructuredLoggingIndicators

	for _, pattern := range strongIndicators {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["structured_logging_detected"] = true
			return
		}
	}
}
