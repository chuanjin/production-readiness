package patterns

// APIGatewayRateLimitPatterns checks for rate limiting in API Gateway configurations
var APIGatewayRateLimitPatterns = []string{
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

// SLOPatterns checks for Service Level Objective configurations
var SLOPatterns = []string{
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

// ErrorBudgetPatterns checks for error budget configurations
var ErrorBudgetPatterns = []string{
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

// SecretsProviderPatterns checks if code uses secrets management services
var SecretsProviderPatterns = []string{
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

// InfraPatterns checks if IaC (Infrastructure as Code) is present
var InfraPatterns = []string{
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

// AWSRegions list of AWS regions
var AWSRegions = []string{
	"us-east-1", "us-east-2", "us-west-1", "us-west-2",
	"af-south-1", "ap-east-1", "ap-south-1", "ap-northeast-3", "ap-northeast-2",
	"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
	"ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1",
	"eu-west-3", "eu-north-1", "me-south-1", "sa-east-1", "us-gov-east-1", "us-gov-west-1",
}

// GCPRegions list of GCP regions
var GCPRegions = []string{
	"us-central1", "us-east1", "us-east4", "us-west1", "us-west2", "us-west3", "us-west4",
	"southamerica-east1", "northamerica-northeast1",
	"europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "europe-north1",
	"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3",
	"asia-southeast1", "asia-southeast2", "australia-southeast1",
}

// AzureRegions list of Azure regions
var AzureRegions = []string{
	"eastus", "eastus2", "southcentralus", "westus2", "westus3", "australiaeast",
	"southeastasia", "northeurope", "westeurope", "uksouth", "ukwest", "francecentral",
	"germanywestcentral", "norwayeast", "switzerlandnorth", "japaneast", "japanwest",
	"centralindia", "southindia", "westindia", "canadacentral", "koreacentral",
}

// ManualStepPatterns indicators of manual deployment steps
var ManualStepPatterns = []string{
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

// MigrationToolPatterns checks for database migration tools
var MigrationToolPatterns = []string{
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

// BackwardCompatPatterns checks for backward compatibility hints
var BackwardCompatPatterns = []string{
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

// MigrationValidationPatterns checks for migration validation steps
var MigrationValidationPatterns = []string{
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

// MutableTags (anti-pattern)
var MutableTags = []string{":latest", ":main", ":master", ":dev", ":develop"}

// VersioningPatterns checks for versioned artifact patterns
var VersioningPatterns = []string{
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

// HealthPatterns checks for health check HTTP endpoints
var HealthPatterns = []string{
	"/health", "\"/health\"", "'/health'",
	"healthcheck", "health-check",
	"endpoint: /health", "path: /health",
	"route('/health')", "get('/health')",
	"@get(\"/health\")", "@route(\"/health\")",
}

// ReadyPatterns checks for readiness endpoints
var ReadyPatterns = []string{
	"/ready", "\"/ready\"", "'/ready'",
	"/readiness", "/readyz",
	"endpoint: /ready", "path: /ready",
	"route('/ready')", "get('/ready')",
	"@get(\"/ready\")", "@route(\"/ready\")",
}

// CorrelationPatterns checks for correlation/trace ID usage
var CorrelationPatterns = []string{
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

// StructuredLoggingPatterns checks for structured logging libraries and patterns
var StructuredLoggingPatterns = []string{
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

// StrongStructuredLoggingIndicators are sufficient on their own
var StrongStructuredLoggingIndicators = []string{
	"structlog", "logrus", "zerolog", "slog", "zap",
	"winston", "pino", "bunyan",
	"serilog", "ecs-logging",
}

// K8sValidKinds for checking resources
var K8sValidKinds = map[string]bool{
	"Pod": true, "Deployment": true, "StatefulSet": true,
	"DaemonSet": true, "Job": true, "CronJob": true,
	"ReplicaSet": true,
}

// NginxIngressRateLimitAnnotations
var NginxIngressRateLimitAnnotations = []string{
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
