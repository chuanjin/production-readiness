package scanner

import (
	"strings"
)

// detectManualSteps checks if documentation contains manual deployment steps
func detectManualSteps(content, relPath string, signals *RepoSignals) {
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

// detectMigrationTool checks for database migration tools
func detectMigrationTool(content, relPath string, signals *RepoSignals) {
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
func detectBackwardCompatibleMigration(content, relPath string, signals *RepoSignals) {
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
func detectMigrationValidation(content, relPath string, signals *RepoSignals) {
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
