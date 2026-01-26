package scanner

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/chuanjin/production-readiness/internal/patterns"
)

// detectAPIGatewayRateLimit checks for rate limiting in API Gateway configurations
func detectAPIGatewayRateLimit(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["api_gateway_rate_limit"] {
		return
	}

	contentLower := strings.ToLower(content)

	// Check for various API Gateway rate limiting patterns
	rateLimitPatterns := patterns.APIGatewayRateLimitPatterns

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
		rateLimitKeys := patterns.RateLimitYAMLKeys
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

	sloPatterns := patterns.SLOPatterns

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
		sloKeys := patterns.SLOYAMLKeys
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

	errorBudgetPatterns := patterns.ErrorBudgetPatterns

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
		budgetKeys := patterns.ErrorBudgetYAMLKeys
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

// detectTimeoutConfiguration checks for timeout configurations in code and config files
func detectTimeoutConfiguration(content, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["timeout_configured"] {
		return
	}

	contentLower := strings.ToLower(content)

	// Check for timeout patterns in code
	timeoutPatterns := patterns.TimeoutPatterns

	for _, pattern := range timeoutPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["timeout_configured"] = true
			return
		}
	}

	// Also check YAML/JSON config files for timeout settings
	ext := strings.ToLower(filepath.Ext(relPath))
	if ext == ExtYAML || ext == ExtYML || ext == ".json" {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		if checkYAMLForTimeout(doc) {
			signals.BoolSignals["timeout_configured"] = true
		}
	}
}

// checkYAMLForTimeout recursively checks YAML for timeout configuration
func checkYAMLForTimeout(obj interface{}) bool {
	switch v := obj.(type) {
	case map[string]interface{}:
		// Check for timeout-related keys
		timeoutKeys := patterns.TimeoutConfigKeys
		for _, key := range timeoutKeys {
			lowerKey := strings.ToLower(key)
			for k := range v {
				if strings.Contains(strings.ToLower(k), lowerKey) {
					return true
				}
			}
		}
		// Recursively check nested objects
		for _, value := range v {
			if checkYAMLForTimeout(value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkYAMLForTimeout(item) {
				return true
			}
		}
	}
	return false
}
