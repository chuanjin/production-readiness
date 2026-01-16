package scanner

import (
	"strings"
)

// detectArtifactVersioning checks for versioned artifact patterns
func detectArtifactVersioning(content, relPath string, signals *RepoSignals) {
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

// detectHealthEndpoints checks for health check HTTP endpoints
func detectHealthEndpoints(content, relPath string, signals *RepoSignals) {
	contentLower := strings.ToLower(content)

	// Detect /health endpoint
	if signals.StringSignals["http_endpoint"] == "" {
		healthPatterns := []string{
			"/health", "\"/health\"", "'/health'",
			"healthcheck", "health-check",
			"endpoint: /health", "path: /health",
			"route('/health')", "get('/health')",
			"@get(\"/health\")", "@route(\"/health\")",
		}

		for _, pattern := range healthPatterns {
			if strings.Contains(contentLower, pattern) {
				signals.StringSignals["http_endpoint"] = "/health"
				break
			}
		}
	}

	// Detect /ready or /readiness endpoint
	if signals.StringSignals["http_endpoint"] == "" {
		readyPatterns := []string{
			"/ready", "\"/ready\"", "'/ready'",
			"/readiness", "/readyz",
			"endpoint: /ready", "path: /ready",
			"route('/ready')", "get('/ready')",
			"@get(\"/ready\")", "@route(\"/ready\")",
		}

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

	correlationPatterns := []string{
		// Common correlation ID names
		"correlation-id", "correlationid", "correlation_id",
		"x-correlation-id", "x-request-id", "x-trace-id",

		// Request ID (similar concept)
		"request-id", "requestid", "request_id",

		// Trace ID (from distributed tracing)
		"trace-id", "traceid", "trace_id", "traceparent",

		// OpenTelemetry
		"opentelemetry", "otel", "trace.traceid",

		// Specific tracing libraries
		"jaeger", "zipkin", "datadog.trace",

		// AWS X-Ray
		"x-amzn-trace-id", "xray",

		// Context propagation
		"propagate", "baggage", "context.context",

		// Logging with correlation
		"logger.with", "log.with", "withfield",
	}

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

	structuredLoggingPatterns := []string{
		// Go libraries
		"logrus", "zap", "zerolog", "slog",

		// Python libraries
		"structlog", "python-json-logger", "pythonjsonlogger",

		// JavaScript/TypeScript
		"winston", "pino", "bunyan",

		// Java libraries
		"logback", "log4j2", "slf4j",

		// .NET libraries
		"serilog", "nlog",

		// Ruby libraries
		"semantic_logger", "ougai",

		// Structured logging patterns
		"log.info", "log.error", "log.warn",
		"logger.info", "logger.error", "logger.warn",
		"withfields", "withfield", "with(", ".with(",

		// JSON logging
		"json.marshal", "json.dumps", "json.stringify",
		"log format: json", "log_format=json", "format=\"json\"",

		// Key-value pairs in logs
		"fields{", "fields:", "attributes{", "context{",

		// ECS (Elastic Common Schema)
		"ecs-logging",
	}

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
	strongIndicators := []string{
		"structlog", "logrus", "zerolog", "slog", "zap",
		"winston", "pino", "bunyan",
		"serilog", "ecs-logging",
	}

	for _, pattern := range strongIndicators {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["structured_logging_detected"] = true
			return
		}
	}
}
