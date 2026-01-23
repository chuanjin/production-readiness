package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRules(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a valid rule file
	validRuleContent := `
id: "test-rule-1"
severity: "high"
category: "security"
title: "Test Rule 1"
description: "This is a test rule"
why_it_matters:
  - "Reason 1"
confidence: "high"
detect:
  all_of:
    - pattern: "test-pattern"
`
	validRulePath := filepath.Join(tempDir, "rule1.yaml")
	if err := os.WriteFile(validRulePath, []byte(validRuleContent), 0o600); err != nil {
		t.Fatalf("Failed to write valid rule file: %v", err)
	}

	// Create a subdirectory with another valid rule
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0o700); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	validRule2Content := `
id: "test-rule-2"
severity: "medium"
`
	validRule2Path := filepath.Join(subDir, "rule2.yml") // .yml extension
	if err := os.WriteFile(validRule2Path, []byte(validRule2Content), 0o600); err != nil {
		t.Fatalf("Failed to write valid rule 2 file: %v", err)
	}

	// Create a non-yaml file (should be ignored)
	ignorePath := filepath.Join(tempDir, "ignore.txt")
	if err := os.WriteFile(ignorePath, []byte("ignored"), 0o600); err != nil {
		t.Fatalf("Failed to write ignore file: %v", err)
	}

	// Create an invalid yaml file
	invalidRuleContent := `
id: "test-rule-3"
severity: "low"
  invalid_indentation: "oops"
`
	invalidRulePath := filepath.Join(tempDir, "invalid.yaml")
	if err := os.WriteFile(invalidRulePath, []byte(invalidRuleContent), 0o600); err != nil {
		t.Fatalf("Failed to write invalid rule file: %v", err)
	}

	// Test case 1: Load rules from directory with mixed files
	// Note: validation of YAML content is done by yaml.Unmarshal, which should return error for invalid YAML
	// Our LoadRules function returns error immediately if any file fails to unmarshal/read.
	// So we expect this to fail because of invalid.yaml
	rules, err := LoadRules(tempDir)
	if err == nil {
		t.Errorf("Expected error due to invalid yaml file, got nil")
	}
	if rules != nil {
		t.Errorf("Expected nil rules on error, got %d rules", len(rules))
	}

	// Test case 2: Load rules from directory with only valid files
	cleanDir := t.TempDir()
	if wErr := os.WriteFile(filepath.Join(cleanDir, "rule1.yaml"), []byte(validRuleContent), 0o600); wErr != nil {
		t.Fatalf("Failed to write clean rule 1: %v", wErr)
	}
	rules, err = LoadRules(cleanDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(rules))
	} else {
		if rules[0].ID != "test-rule-1" {
			t.Errorf("Expected rule ID 'test-rule-1', got '%s'", rules[0].ID)
		}
		if rules[0].Severity != High { // Assuming High const is available
			t.Errorf("Expected severity 'high', got '%s'", rules[0].Severity)
		}
	}

	// Test case 3: Recursive loading
	// Already tested partially, but let's be explicit
	recursiveDir := t.TempDir()
	if mkErr := os.MkdirAll(filepath.Join(recursiveDir, "a", "b"), 0o700); mkErr != nil {
		t.Fatal(mkErr)
	}
	if wErr := os.WriteFile(filepath.Join(recursiveDir, "a", "rule.yaml"), []byte(validRuleContent), 0o600); wErr != nil {
		t.Fatal(wErr)
	}
	if wErr := os.WriteFile(filepath.Join(recursiveDir, "a", "b", "rule.yml"), []byte(validRule2Content), 0o600); wErr != nil {
		t.Fatal(wErr)
	}
	rules, err = LoadRules(recursiveDir)
	if err != nil {
		t.Fatalf("Expected no error for recursive load, got %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(rules))
	}

	// Test case 4: Non-existent directory
	_, err = LoadRules("/path/to/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}
