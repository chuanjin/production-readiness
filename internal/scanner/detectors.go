package scanner

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	// Auto-register all detectors
	// Add new detectors here - just one line per detector!
	registerDetector(detectSecretsProvider)
	registerDetector(detectInfrastructure)
	registerDetector(detectRegions)
	registerDetector(detectManualSteps)
	registerDetector(detectK8sDeploymentStrategy)
	registerDetector(detectArtifactVersioning)
	registerDetector(detectHealthEndpoints)
	registerDetector(detectK8sProbes)
	registerDetector(detectCorrelationId)
	registerDetector(detectStructuredLogging)
	registerDetector(detectIngressRateLimit)
	registerDetector(detectAPIGatewayRateLimit)
	registerDetector(detectSLOConfig)
	registerDetector(detectErrorBudget)
	registerDetector(detectMigrationTool)
	registerDetector(detectBackwardCompatibleMigration)
	registerDetector(detectMigrationValidation)
}

// detectSecretsProvider checks if code uses secrets management services
func detectSecretsProvider(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["secrets_provider_detected"] {
		return
	}

	secretsProviderPatterns := []string{
		// AWS Secrets Manager
		"aws-sdk", "aws/secretsmanager", "GetSecretValue", "secretsmanager",
		"AWS::SecretsManager",

		// HashiCorp Vault
		"hashicorp/vault", "vault.NewClient", "vault/api",

		// Google Secret Manager
		"cloud.google.com/go/secretmanager", "secretmanager.NewClient",
		"google-cloud/secret-manager",

		// Azure Key Vault
		"azure-keyvault", "azure/keyvault", "KeyVaultClient",

		// Doppler
		"doppler.com", "DopplerSDK", "@dopplerhq",

		// Infisical
		"infisical", "infisical-sdk",

		// 1Password
		"1password", "op://",

		// Generic secrets management
		"sealed-secrets", "external-secrets", "secrets-store-csi",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range secretsProviderPatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			signals.BoolSignals["secrets_provider_detected"] = true
			return
		}
	}
}

// detectK8sDeploymentStrategy checks Kubernetes deployment files for strategy
func detectK8sDeploymentStrategy(content string, relPath string, signals *RepoSignals) {
	if signals.StringSignals["k8s_deployment_strategy"] != "" {
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

// detectArtifactVersioning checks for versioned artifact patterns
func detectArtifactVersioning(content string, relPath string, signals *RepoSignals) {
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

// detectInfrastructure checks if IaC (Infrastructure as Code) is present
func detectInfrastructure(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["infra_as_code_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	infraPatterns := []string{
		// Terraform
		"terraform", "provider \"", "resource \"", "module \"",

		// CloudFormation
		"aws::cloudformation", "awscloudformation", "resources:",

		// Pulumi
		"pulumi", "@pulumi/",

		// CDK
		"aws-cdk", "@aws-cdk/",

		// Kubernetes/Helm
		"apiversion:", "kind: deployment", "kind: service",

		// Ansible
		"ansible", "playbook",
	}

	for _, pattern := range infraPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["infra_as_code_detected"] = true
			return
		}
	}
}

// detectRegions counts the number of unique cloud regions configured
func detectRegions(content string, relPath string, signals *RepoSignals) {
	regions := make(map[string]bool)

	contentLower := strings.ToLower(content)

	// AWS regions
	awsRegions := []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-north-1",
		"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ap-northeast-2",
		"sa-east-1", "ca-central-1", "ap-south-1",
	}

	// GCP regions
	gcpRegions := []string{
		"us-central1", "us-east1", "us-west1", "us-east4",
		"europe-west1", "europe-west2", "europe-west3", "europe-north1",
		"asia-east1", "asia-southeast1", "asia-northeast1",
	}

	// Azure regions
	azureRegions := []string{
		"eastus", "eastus2", "westus", "westus2", "centralus",
		"northeurope", "westeurope", "uksouth", "ukwest",
		"southeastasia", "eastasia", "japaneast", "japanwest",
	}

	allRegions := append(append(awsRegions, gcpRegions...), azureRegions...)

	for _, region := range allRegions {
		if strings.Contains(contentLower, region) {
			regions[region] = true
		}
	}

	// Update the count (only if we found more regions than before)
	if len(regions) > signals.IntSignals["region_count"] {
		signals.IntSignals["region_count"] = len(regions)
	}
}

// detectManualSteps checks if documentation contains manual deployment steps
func detectManualSteps(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["manual_steps_documented"] {
		return
	}

	// Only check documentation files
	fileName := strings.ToLower(relPath)
	isDocFile := strings.Contains(fileName, "readme") ||
		strings.Contains(fileName, "doc") ||
		strings.Contains(fileName, "deploy") ||
		strings.Contains(fileName, "setup") ||
		strings.Contains(fileName, "install") ||
		strings.HasSuffix(fileName, ".md") ||
		strings.HasSuffix(fileName, ".txt")

	if !isDocFile {
		return
	}

	contentLower := strings.ToLower(content)

	// Patterns indicating manual steps
	manualStepPatterns := []string{
		// Step-by-step instructions
		"step 1", "step 2", "1.", "2.", "3.",
		"first,", "then,", "next,", "finally,",

		// Manual actions
		"manually", "by hand", "login to", "navigate to",
		"click on", "open the", "go to the console",
		"ssh into", "copy the file", "run this command",

		// Console/UI instructions
		"in the console", "in the dashboard", "in the ui",
		"from the web interface", "using the portal",

		// Manual verification
		"verify that", "check that", "make sure",
		"confirm that", "ensure that",

		// Manual configuration
		"edit the file", "update the", "change the",
		"set the value", "configure manually",
	}

	// Count matches to avoid false positives
	matches := 0
	for _, pattern := range manualStepPatterns {
		if strings.Contains(contentLower, pattern) {
			matches++
			if matches >= 3 { // Need at least 3 indicators
				signals.BoolSignals["manual_steps_documented"] = true
				return
			}
		}
	}
}

// detectHealthEndpoints checks for health check HTTP endpoints
func detectHealthEndpoints(content string, relPath string, signals *RepoSignals) {
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

// detectK8sProbes checks for Kubernetes liveness/readiness probes
func detectK8sProbes(content string, relPath string, signals *RepoSignals) {
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
	validKinds := map[string]bool{
		"Pod": true, "Deployment": true, "StatefulSet": true,
		"DaemonSet": true, "Job": true, "CronJob": true,
		"ReplicaSet": true,
	}

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

// detectCorrelationId checks for correlation/trace ID usage
func detectCorrelationId(content string, relPath string, signals *RepoSignals) {
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
func detectStructuredLogging(content string, relPath string, signals *RepoSignals) {
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

// detectIngressRateLimit checks for rate limiting in Kubernetes Ingress
func detectIngressRateLimit(content string, relPath string, signals *RepoSignals) {
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
			rateLimitAnnotations := []string{
				"nginx.ingress.kubernetes.io/limit-rps",
				"nginx.ingress.kubernetes.io/limit-rpm",
				"nginx.ingress.kubernetes.io/limit-connections",
				"nginx.ingress.kubernetes.io/limit-burst-multiplier",

				// Traefik rate limiting
				"traefik.ingress.kubernetes.io/rate-limit",

				// Kong rate limiting
				"konghq.com/plugins",
				"rate-limiting.plugin.konghq.com",
			}

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

// detectAPIGatewayRateLimit checks for rate limiting in API Gateway configurations
func detectAPIGatewayRateLimit(content string, relPath string, signals *RepoSignals) {
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
	if ext == ".yaml" || ext == ".yml" {
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
func detectSLOConfig(content string, relPath string, signals *RepoSignals) {
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
	if ext == ".yaml" || ext == ".yml" {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		// Check for OpenSLO format
		if kind, ok := doc["kind"].(string); ok {
			if strings.ToLower(kind) == "slo" || strings.ToLower(kind) == "servicelevelobjective" {
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

// detectErrorBudget checks for error budget configurations
func detectErrorBudget(content string, relPath string, signals *RepoSignals) {
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
	if ext == ".yaml" || ext == ".yml" {
		var doc map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return
		}

		if checkYAMLForErrorBudget(doc) {
			signals.BoolSignals["error_budget_detected"] = true
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

// detectMigrationTool checks for database migration tools
func detectMigrationTool(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["migration_tool_detected"] {
		return
	}

	contentLower := strings.ToLower(content)

	migrationToolPatterns := []string{
		// Go migration tools
		"golang-migrate", "migrate.up", "migrate.down",
		"goose", "sql-migrate",

		// Node.js/TypeScript
		"knex", "sequelize", "typeorm", "prisma migrate",
		"db-migrate", "umzug",

		// Python
		"alembic", "django.db.migrations", "flask-migrate",
		"yoyo-migrations", "sqlalchemy-migrate",

		// Ruby
		"activerecord::migration", "rake db:migrate",

		// Java
		"flyway", "liquibase",

		// .NET
		"entity framework", "fluentmigrator",

		// Generic patterns
		"migrations/", "migration.sql", "schema_migrations",
		"up.sql", "down.sql", "migrate up", "migrate down",
		"create_table", "alter_table", "add_column", "drop_column",
	}

	for _, pattern := range migrationToolPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.BoolSignals["migration_tool_detected"] = true
			return
		}
	}
}

// detectBackwardCompatibleMigration checks for backward compatibility hints
func detectBackwardCompatibleMigration(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["backward_compatible_migration_hint"] {
		return
	}

	contentLower := strings.ToLower(content)

	backwardCompatPatterns := []string{
		// Explicit backward compatibility
		"backward compatible", "backwards compatible",
		"backward-compatible", "backwards-compatible",
		"zero-downtime", "zero downtime",

		// Expand-contract pattern
		"expand and contract", "expand-contract",
		"dual-write", "dual write", "shadow write",

		// Safe migration practices
		"nullable", "null: true", "default:", "default value",
		// We treat these as weaker indicators that must appear in combination
		// (e.g. ADD COLUMN + NULL/DEFAULT) to avoid false positives.
		"add column", "null", "default",

		// Incremental changes
		"incremental migration", "phased migration",
		"blue-green", "canary",

		// Documentation about compatibility
		"safe to deploy", "rollback safe", "reversible",
		"no breaking change", "non-breaking",

		// Feature flags for migrations
		"feature flag", "feature toggle", "flag:",
	}

	matchCount := 0
	for _, pattern := range backwardCompatPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Strong indicators
			if strings.Contains(pattern, "backward") ||
				strings.Contains(pattern, "zero-downtime") ||
				strings.Contains(pattern, "expand-contract") {
				signals.BoolSignals["backward_compatible_migration_hint"] = true
				return
			}
			// Weaker indicators - need multiple
			if matchCount >= 2 {
				signals.BoolSignals["backward_compatible_migration_hint"] = true
				return
			}
		}
	}
}

// detectMigrationValidation checks for migration validation steps
func detectMigrationValidation(content string, relPath string, signals *RepoSignals) {
	if signals.BoolSignals["migration_validation_step"] {
		return
	}

	contentLower := strings.ToLower(content)

	validationPatterns := []string{
		// Explicit validation
		"validate", "validation", "verify migration",
		"check migration", "test migration",

		// Dry run
		"dry-run", "dry run", "--dry-run", "dryrun",
		"simulate", "plan", "preview",

		// Migration testing
		"migration test", "test:migration", "migration_test",
		"test_migration", "test_migration_validation",

		// Rollback testing
		"rollback test", "test rollback", "rollback", "revert",
		"migration down", "migrate down",

		// Data validation
		"data integrity", "consistency check", "validate data",
		"check constraint", "foreign key check",

		// Schema validation
		"schema validation", "validate schema", "schema check",

		// CI/CD validation
		"migration ci", "ci migration", "test:db",

		// Safety checks
		"pre-migration", "post-migration", "migration hook",
		"before_migrate", "after_migrate",

		// Backup before migration
		"backup before", "snapshot before", "dump before",
	}

	matchCount := 0
	for _, pattern := range validationPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Strong indicators
			if strings.Contains(pattern, "validate") ||
				strings.Contains(pattern, "test") ||
				strings.Contains(pattern, "dry-run") ||
				strings.Contains(pattern, "rollback") {
				signals.BoolSignals["migration_validation_step"] = true
				return
			}
			// Weaker indicators - need multiple
			if matchCount >= 2 {
				signals.BoolSignals["migration_validation_step"] = true
				return
			}
		}
	}
}
