package scanner

import (
	"testing"
)

func TestDetectK8sDeploymentStrategy(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		relPath  string
		expected string
	}{
		{
			name: "RollingUpdate strategy",
			content: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  strategy:
    type: RollingUpdate
`,
			relPath:  "deploy.yaml",
			expected: "RollingUpdate",
		},
		{
			name: "Recreate strategy",
			content: `
apiVersion: apps/v1
kind: Deployment
spec:
  strategy:
    type: Recreate
`,
			relPath:  "deploy.yml",
			expected: "Recreate",
		},
		{
			name: "Not a deployment",
			content: `
kind: Service
spec:
  type: ClusterIP
`,
			relPath:  "service.yaml",
			expected: "",
		},
		{
			name: "Not a YAML file",
			content: `
strategy: RollingUpdate
`,
			relPath:  "readme.md",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				StringSignals: make(map[string]string),
			}
			detectK8sDeploymentStrategy(tt.content, tt.relPath, signals)

			if got := signals.StringSignals["k8s_deployment_strategy"]; got != tt.expected {
				t.Errorf("detectK8sDeploymentStrategy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectIngressRateLimit(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		relPath  string
		expected bool
	}{
		{
			name: "Nginx rate limit annotation",
			content: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/limit-rps: "10"
`,
			relPath:  "ingress.yaml",
			expected: true,
		},
		{
			name: "Kong rate limit plugin",
			content: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    konghq.com/plugins: "my-rate-limit"
`,
			relPath:  "ingress.yml",
			expected: true,
		},
		{
			name: "No rate limit",
			content: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: simple-ingress
`,
			relPath:  "ingress.yaml",
			expected: false,
		},
		{
			name: "Not an Ingress",
			content: `
kind: Service
metadata:
  annotations:
    nginx.ingress.kubernetes.io/limit-rps: "10"
`,
			relPath:  "service.yaml",
			expected: false,
		},
		{
			name: "Invalid YAML content",
			content: `
invalid: yaml: [[[
`,
			relPath:  "ingress.yaml",
			expected: false,
		},
		{
			name: "Ingress without metadata",
			content: `
apiVersion: networking.k8s.io/v1
kind: Ingress
spec:
  rules: []
`,
			relPath:  "ingress.yaml",
			expected: false,
		},
		{
			name: "Ingress without annotations",
			content: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: simple-ingress
spec:
  rules: []
`,
			relPath:  "ingress.yaml",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectIngressRateLimit(tt.content, tt.relPath, signals)

			if got := signals.GetBool("ingress_rate_limit"); got != tt.expected {
				t.Errorf("detectIngressRateLimit() = %v, want %v", got, tt.expected)
			}
		})
	}

	// Test early return when signal already detected
	t.Run("Early return when already detected", func(t *testing.T) {
		signals := &RepoSignals{
			BoolSignals: map[string]bool{"ingress_rate_limit": true},
		}
		content := `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/limit-rps: "10"
`
		detectIngressRateLimit(content, "ingress.yaml", signals)
		// Should still be true, function returns early
		if !signals.GetBool("ingress_rate_limit") {
			t.Error("expected signal to remain true")
		}
	})
}

func TestDetectK8sProbes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Liveness probe detected",
			content: `
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        livenessProbe:
          httpGet:
            path: /health
`,
			expected: true,
		},
		{
			name: "Readiness probe detected",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    readinessProbe:
      httpGet:
        path: /ready
`,
			expected: true,
		},
		{
			name: "No probes",
			content: `
apiVersion: v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectK8sProbes(tt.content, "deployment.yaml", signals)

			if signals.GetBool("k8s_probe_defined") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("k8s_probe_defined"))
			}
		})
	}
}

func TestDetectResourceLimits(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Pod with limits detected",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    resources:
      limits:
        memory: "128Mi"
`,
			expected: true,
		},
		{
			name: "Deployment with limits detected",
			content: `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        resources:
          limits:
            cpu: "500m"
`,
			expected: true,
		},
		{
			name: "No limits defined",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    resources:
      requests:
        cpu: "100m"
`,
			expected: false,
		},
		{
			name: "No resources defined",
			content: `
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectResourceLimits(tt.content, "deploy.yaml", signals)

			if signals.GetBool("k8s_resource_limits_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("k8s_resource_limits_detected"))
			}
		})
	}
}
