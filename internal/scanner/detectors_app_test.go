package scanner

import (
	"testing"
)

func TestDetectVersionedArtifacts(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Semantic version detected",
			content:  `image: myapp:v1.2.3`,
			expected: true,
		},
		{
			name:     "SHA256 digest detected",
			content:  `image: myapp@sha256:abc123`,
			expected: true,
		},
		{
			name:     "Latest tag - not versioned",
			content:  `image: myapp:latest`,
			expected: false,
		},
		{
			name:     "Git tag detected",
			content:  `git tag v1.0.0`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectArtifactVersioning(tt.content, "deploy.yaml", signals)

			if signals.GetBool("versioned_artifacts") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("versioned_artifacts"))
			}
		})
	}
}

func TestDetectHealthEndpoints(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedEndpoint string
	}{
		{
			name:             "Health endpoint detected",
			content:          `app.get('/health', handler)`,
			expectedEndpoint: "/health",
		},
		{
			name:             "Ready endpoint detected",
			content:          `@route('/ready')\ndef ready():`,
			expectedEndpoint: "/ready",
		},
		{
			name:             "No endpoint",
			content:          `app.get('/api/users', handler)`,
			expectedEndpoint: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				StringSignals: make(map[string]string),
			}
			detectHealthEndpoints(tt.content, "server.js", signals)

			if signals.GetString("http_endpoint") != tt.expectedEndpoint {
				t.Errorf("expected %q, got %q", tt.expectedEndpoint, signals.GetString("http_endpoint"))
			}
		})
	}
}

func TestDetectCorrelationId(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "X-Request-ID header",
			content:  `req.Header.Get("X-Request-ID")`,
			expected: true,
		},
		{
			name:     "Correlation ID variable",
			content:  `correlationId := generateID()`,
			expected: true,
		},
		{
			name:     "OpenTelemetry trace",
			content:  `import "go.opentelemetry.io/otel"`,
			expected: true,
		},
		{
			name:     "No correlation ID",
			content:  `log.Println("hello")`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectCorrelationID(tt.content, "handler.go", signals)

			if signals.GetBool("correlation_id_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("correlation_id_detected"))
			}
		})
	}
}

func TestDetectStructuredLogging(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Logrus detected",
			content:  `import "github.com/sirupsen/logrus"`,
			expected: true,
		},
		{
			name:     "Winston detected",
			content:  `const winston = require('winston');`,
			expected: true,
		},
		{
			name:     "Plain console.log",
			content:  `console.log("hello");`,
			expected: false,
		},
		{
			name:     "Zap logger",
			content:  `logger, _ := zap.NewProduction()`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectStructuredLogging(tt.content, "logger.go", signals)

			if signals.GetBool("structured_logging_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("structured_logging_detected"))
			}
		})
	}
}
