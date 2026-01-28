package scanner

import (
	"strings"

	"github.com/chuanjin/production-readiness/internal/patterns"
)

// detectManualSteps checks if documentation contains manual deployment steps
func detectManualSteps(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("manual_steps_documented") {
		return
	}

	// Only check documentation files
	fileName := strings.ToLower(relPath)

	isDocFile := false
	for _, keyword := range patterns.DocFileKeywords {
		if strings.Contains(fileName, keyword) {
			isDocFile = true
			break
		}
	}

	if !isDocFile {
		for _, ext := range patterns.DocFileExtensions {
			if strings.HasSuffix(fileName, ext) {
				isDocFile = true
				break
			}
		}
	}

	if !isDocFile {
		return
	}

	contentLower := strings.ToLower(content)

	// Patterns indicating manual steps
	manualStepPatterns := patterns.ManualStepPatterns

	// Count matches to avoid false positives
	matches := 0
	for _, pattern := range manualStepPatterns {
		if strings.Contains(contentLower, pattern) {
			matches++
			if matches >= 3 { // Need at least 3 indicators
				signals.SetBool("manual_steps_documented", true)
				return
			}
		}
	}
}

// detectMigrationTool checks for database migration tools
func detectMigrationTool(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("migration_tool_detected") {
		return
	}

	contentLower := strings.ToLower(content)

	migrationToolPatterns := patterns.MigrationToolPatterns

	for _, pattern := range migrationToolPatterns {
		if strings.Contains(contentLower, pattern) {
			signals.SetBool("migration_tool_detected", true)
			return
		}
	}
}

// detectBackwardCompatibleMigration checks for backward compatibility hints
func detectBackwardCompatibleMigration(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("backward_compatible_migration_hint") {
		return
	}

	contentLower := strings.ToLower(content)

	backwardCompatPatterns := patterns.BackwardCompatPatterns

	matchCount := 0
	for _, pattern := range backwardCompatPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Strong indicators
			if strings.Contains(pattern, "backward") ||
				strings.Contains(pattern, "zero-downtime") ||
				strings.Contains(pattern, "expand-contract") {
				signals.SetBool("backward_compatible_migration_hint", true)
				return
			}
			// Weaker indicators - need multiple
			if matchCount >= 2 {
				signals.SetBool("backward_compatible_migration_hint", true)
				return
			}
		}
	}
}

// detectMigrationValidation checks for migration validation steps
func detectMigrationValidation(content, relPath string, signals *RepoSignals) {
	if signals.GetBool("migration_validation_step") {
		return
	}

	contentLower := strings.ToLower(content)

	validationPatterns := patterns.MigrationValidationPatterns

	matchCount := 0
	for _, pattern := range validationPatterns {
		if strings.Contains(contentLower, pattern) {
			matchCount++
			// Strong indicators
			if strings.Contains(pattern, "validate") ||
				strings.Contains(pattern, "test") ||
				strings.Contains(pattern, "dry-run") ||
				strings.Contains(pattern, "rollback") {
				signals.SetBool("migration_validation_step", true)
				return
			}
			// Weaker indicators - need multiple
			if matchCount >= 2 {
				signals.SetBool("migration_validation_step", true)
				return
			}
		}
	}
}
