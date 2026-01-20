package scanner

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/chuanjin/production-readiness/internal/patterns"
)

// detectK8sDeploymentStrategy checks Kubernetes deployment files for strategy
func detectK8sDeploymentStrategy(content, relPath string, signals *RepoSignals) {
	if signals.StringSignals["k8s_deployment_strategy"] != "" {
		return
	}

	// Only check YAML files
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext != ExtYAML && ext != ExtYML {
		return
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return
	}

	kind, ok := doc["kind"].(string)
	if !ok || kind != "Deployment" {
		return
	}

	if spec, ok := doc["spec"].(map[string]interface{}); ok {
		if strategy, ok := spec["strategy"].(map[string]interface{}); ok {
			if strategyType, ok := strategy["type"].(string); ok {
				signals.StringSignals["k8s_deployment_strategy"] = strategyType
			}
		}
	}
}

// detectK8sProbes checks for Kubernetes liveness/readiness probes
func detectK8sProbes(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["k8s_probe_defined"] {
		return
	}

	// Only check YAML files
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext != ".yaml" && ext != ".yml" {
		return
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return
	}

	// Check if it's a Kubernetes resource with containers
	kind, ok := doc["kind"].(string)
	if !ok {
		return
	}

	// Look for probes in Pod, Deployment, StatefulSet, DaemonSet, etc.
	validKinds := patterns.K8sValidKinds

	if !validKinds[kind] {
		return
	}

	// Navigate to containers
	var containers []interface{}

	if spec, ok := doc["spec"].(map[string]interface{}); ok {
		// For Deployments, StatefulSets, etc., probes are in spec.template.spec.containers
		if template, ok := spec["template"].(map[string]interface{}); ok {
			if templateSpec, ok := template["spec"].(map[string]interface{}); ok {
				if c, ok := templateSpec["containers"].([]interface{}); ok {
					containers = c
				}
			}
		} else if c, ok := spec["containers"].([]interface{}); ok {
			// For Pods, probes are directly in spec.containers
			containers = c
		}
	}

	// Check if any container has probes
	for _, container := range containers {
		if c, ok := container.(map[string]interface{}); ok {
			// Check for livenessProbe or readinessProbe
			if _, hasLiveness := c["livenessProbe"]; hasLiveness {
				signals.BoolSignals["k8s_probe_defined"] = true
				return
			}
			if _, hasReadiness := c["readinessProbe"]; hasReadiness {
				signals.BoolSignals["k8s_probe_defined"] = true
				return
			}
		}
	}
}

// detectIngressRateLimit checks for rate limiting in Kubernetes Ingress
func detectIngressRateLimit(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["ingress_rate_limit"] {
		return
	}

	// Only check YAML files
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext != ".yaml" && ext != ".yml" {
		return
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return
	}

	// Check if it's an Ingress resource
	kind, ok := doc["kind"].(string)
	if !ok || kind != "Ingress" {
		return
	}

	// Check annotations for rate limiting
	if metadata, ok := doc["metadata"].(map[string]interface{}); ok {
		if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
			// NGINX Ingress rate limiting annotations
			rateLimitAnnotations := patterns.NginxIngressRateLimitAnnotations

			for _, annotation := range rateLimitAnnotations {
				if _, exists := annotations[annotation]; exists {
					signals.BoolSignals["ingress_rate_limit"] = true
					return
				}
			}

			// Also check if Kong plugins annotation contains rate-limiting
			if plugins, ok := annotations["konghq.com/plugins"].(string); ok {
				if strings.Contains(strings.ToLower(plugins), "rate-limit") {
					signals.BoolSignals["ingress_rate_limit"] = true
					return
				}
			}
		}
	}
}

// detectResourceLimits checks for Kubernetes resource limits configurations
func detectResourceLimits(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["k8s_resource_limits_detected"] {
		return
	}

	ext := strings.ToLower(filepath.Ext(relPath))
	if ext != ExtYAML && ext != ExtYML {
		return
	}

	var doc map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return
	}

	// Check if it's a Kubernetes resource
	kind, ok := doc["kind"].(string)
	if !ok {
		return
	}

	validKinds := patterns.K8sValidKinds

	if !validKinds[kind] {
		return
	}

	// Navigate to containers
	var containers []interface{}

	if spec, ok := doc["spec"].(map[string]interface{}); ok {
		if template, ok := spec["template"].(map[string]interface{}); ok {
			if templateSpec, ok := template["spec"].(map[string]interface{}); ok {
				if c, ok := templateSpec["containers"].([]interface{}); ok {
					containers = c
				}
			}
		} else if c, ok := spec["containers"].([]interface{}); ok {
			// For Pods, probes are directly in spec.containers
			containers = c
		}
	}

	for _, container := range containers {
		if c, ok := container.(map[string]interface{}); ok {
			if resources, ok := c["resources"].(map[string]interface{}); ok {
				if limits, ok := resources["limits"].(map[string]interface{}); ok {
					// Check if cpu or memory limits are defined
					_, hasCPU := limits["cpu"]
					_, hasMemory := limits["memory"]
					if hasCPU || hasMemory {
						signals.BoolSignals["k8s_resource_limits_detected"] = true
						return
					}
				}
			}
		}
	}
}
