package scanner

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// detectAPIGatewayRateLimit checks for rate limiting in API Gateway configurations
func detectAPIGatewayRateLimit(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["api_gateway_rate_limit"] {
		return
	}

	contentLower := strings.ToLower(content)

	// Check for various API Gateway rate limiting patterns
	rateLimitPatterns := []string{
		// AWS API Gateway
		"throttlesettings", "throttle", "ratelimit", "burstlimit",
		"aws::apigateway", "usage plan", "usageplan",

		// Kong
		"rate-limiting", "rate_limiting", "kong-plugin-rate-limiting",

		// Express (Node.js)
		"express-rate-limit", "rate-limiter", "ratelimit(",

		// Go libraries
		"golang.org/x/time/rate", "rate.limiter", "ratelimit.new",
		"throttled", "tollbooth",

		// Python libraries
		"flask-limiter", "django-ratelimit", "slowapi",

		// Redis rate limiting
		"redis-rate-limit", "redis:incr", "redis.incr",

		// NGINX rate limiting
		"limit_req", "limit_conn", "limit_rate",

		// Envoy rate limiting
		"envoy.filters.http.ratelimit", "rate_limit_service",

		// Cloud provider rate limiting
		"cloudfront.ratelimit", "azure.ratelimit",

		// Generic patterns
		"requests per second", "requests per minute",
		"max_requests", "rate_limit", "throttle_rate",
	}

	for _, pattern := range rateLimitPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["api_gateway_rate_limit"] = true
			return
		}
	}

	// Also check YAML for API Gateway configs
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext == ExtYAML || ext == ExtYML {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		// Check for rate limit in various gateway configs
		if checkYAMLForRateLimit(doc) {
			signals.BoolSignals["api_gateway_rate_limit"] = true
		}
	}
}

// checkYAMLForRateLimit recursively checks YAML for rate limiting config
func checkYAMLForRateLimit(obj interface{}) bool {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Check for rate limit keys
		rateLimitKeys := []string{
			"ratelimit", "rate_limit", "rate-limit",
			"throttle", "quota", "limit",
		}
		for _, key := range rateLimitKeys {
			if _, exists := v[key]; exists {
				return true
			}
		}
		// Recursively check nested objects
		for _, value := range v {
			if checkYAMLForRateLimit(value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkYAMLForRateLimit(item) {
				return true
			}
		}
	}
	return false
}

// detectSLOConfig checks for Service Level Objective configurations
func detectSLOConfig(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["slo_config_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	sloPatterns := []string{
		// SLO/SLI keywords
		"slo:", "sli:", "service level objective", "service level indicator",
		"slo_config", "slo-config", "sloconfig",

		// OpenSLO format
		"openslo", "kind: slo", "apiversion: openslo",

		// Prometheus-based SLO
		"sloth", "pyrra", "slo-libsonnet",

		// Cloud provider SLO
		"google_monitoring_slo", "aws_servicelevelobjective",
		"azurerm_monitor_slo",

		// SLO metrics
		"availability_slo", "latency_slo", "error_rate_slo",
		"slo_target", "slo_threshold", "objective:",

		// SLO tools
		"nobl9", "lightstep", "datadog slo",

		// Common SLO patterns
		"99.9%", "99.95%", "99.99%", "four nines", "three nines",
		"uptime_target", "availability_target",
	}

	matchCount := 0
	for _, pattern := range sloPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Strong indicators - single match is enough
			if strings.Contains(pattern, "slo") || strings.Contains(pattern, "objective") {
				signals.BoolSignals["slo_config_detected"] = true
				return
			}
			// Weak indicators - need multiple matches
			if matchCount >= 2 {
				signals.BoolSignals["slo_config_detected"] = true
				return
			}
		}
	}

	// Also check YAML structure for SLO configs
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext == ExtYAML || ext == ExtYML {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		// Check for OpenSLO format
		if kind, ok := doc["kind"].(string); ok {
			if strings.EqualFold(kind, "slo") || strings.EqualFold(kind, "servicelevelobjective") {
				signals.BoolSignals["slo_config_detected"] = true
				return
			}
		}

		// Check for SLO-related keys
		if checkYAMLForSLO(doc) {
			signals.BoolSignals["slo_config_detected"] = true
		}
	}
}

// checkYAMLForSLO recursively checks YAML for SLO configuration
func checkYAMLForSLO(obj interface{}) bool {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Check for SLO-related keys
		sloKeys := []string{
			"slo", "sli", "objective", "objectives",
			"service_level_objective", "service_level_indicator",
			"target", "availability", "latency_target",
		}
		for _, key := range sloKeys {
			lowerKey := strings.ToLower(key)
			for k := range v {
				if strings.Contains(strings.ToLower(k), lowerKey) {
					return true
				}
			}
		}
		// Recursively check nested objects
		for _, value := range v {
			if checkYAMLForSLO(value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkYAMLForSLO(item) {
				return true
			}
		}
	}
	return false
}

// detectErrorBudget checks for error budget configurations
func detectErrorBudget(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["error_budget_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	errorBudgetPatterns := []string{
		// Error budget keywords
		"error_budget", "error-budget", "errorbudget",
		"error budget", "budget:",

		// Error budget policies
		"error_budget_policy", "budget_policy", "burn_rate",
		"burnrate", "burn-rate",

		// Error budget calculation
		"remaining_budget", "budget_remaining", "budget_spent",
		"budget_consumption", "error_rate_threshold",

		// Alerting based on error budget
		"error_budget_alert", "budget_exhausted", "budget_burn",

		// SRE tools with error budgets
		"sloth", "pyrra", "nobl9", "openslo",

		// Prometheus error budget queries
		"error_budget{", "slo_error_budget",

		// Cloud provider error budgets
		"google_monitoring_slo", "consumed_budget",
	}

	for _, pattern := range errorBudgetPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["error_budget_detected"] = true
			return
		}
	}

	// Also check YAML structure for error budget configs
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext == ExtYAML || ext == ExtYML {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		if checkYAMLForErrorBudget(doc) {
			signals.BoolSignals["error_budget_detected"] = true
		}
	}
}

// checkYAMLForErrorBudget recursively checks YAML for error budget configuration
func checkYAMLForErrorBudget(obj interface{}) bool {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Check for error budget keys
		budgetKeys := []string{
			"error_budget", "errorbudget", "error-budget",
			"budget", "burn_rate", "burnrate",
		}
		for _, key := range budgetKeys {
			lowerKey := strings.ToLower(key)
			for k := range v {
				if strings.Contains(strings.ToLower(k), lowerKey) {
					return true
				}
			}
		}
		// Recursively check nested objects
		for _, value := range v {
			if checkYAMLForErrorBudget(value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkYAMLForErrorBudget(item) {
				return true
			}
		}
	}
	return false
}
