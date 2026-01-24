package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanCmd(t *testing.T) {
	// Setup temporary workspace
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Clean up: restore working directory
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatalf("Failed to restore working directory: %v", err)
		}
	}()

	// Change to temp dir so "rules" lookup works
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create dummy rules directory
	if err := os.Mkdir("rules", 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy rule
	ruleContent := `
id: "test-rule"
severity: "high"
category: "reliability"
title: "Test Rule"
description: "Detects test pattern"
why_it_matters:
  - "Testing purposes"
confidence: "high"
detect:
  all_of:
    - pattern: "TEST_PATTERN"
`
	if err := os.WriteFile(filepath.Join("rules", "test-rule.yaml"), []byte(ruleContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create dummy code to scan
	// Use string literal to ensure it is picked up if comments are ignored
	if err := os.WriteFile("main.go", []byte("package main\n\nvar _ = \"TEST_PATTERN\""), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name             string
		args             []string
		expectedContains []string // List of strings that must appear in output
		wantErr          bool
	}{
		{
			name: "Scan markdown output",
			args: []string{"root"},
			// We check for "Total: 1 rules" to confirm rule was loaded and processed.
			// whether it triggers or passes depends on scanner exact logic, but "Total: 1" confirms integration.
			expectedContains: []string{
				"Production Readiness Report",
				"Total: 1 rules",
			},
		},
		{
			name: "Scan json output",
			args: []string{"root", "--format", "json"},
			// JSON should always contain the rule definition/result in some form
			expectedContains: []string{
				"\"id\": \"test-rule\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			// We must pass "scan" and then args.
			// args in test case includes path "root" (which we need to handle? No, we are in tempDir)
			// Wait, the "scan" command takes [path] argument.
			// If we pass ".", it scans current dir.
			// My test logic uses "root"? No, the args in struct above inputs "root"?
			// Ah, I used "root" as placeholder in my thought?
			// Let's use "." as path.

			realArgs := []string{"scan", "."}
			// Append flags
			if len(tt.args) > 1 { // Assuming first arg is path
				realArgs = append(realArgs, tt.args[1:]...)
			}

			rootCmd.SetArgs(realArgs)

			if err := rootCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("rootCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			output := buf.String()
			for _, exp := range tt.expectedContains {
				if !strings.Contains(output, exp) {
					t.Errorf("Output missing expected content '%s'.\nGot:\n%s", exp, output)
				}
			}
		})
	}
}
