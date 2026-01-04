package rules

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadRules(rulesDir string) ([]Rule, error) {
	var rules []Rule

	// Clean and validate rulesDir to prevent directory traversal
	rulesDir, err := filepath.Abs(rulesDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Validate path is within rulesDir to prevent traversal
		if !strings.HasPrefix(path, rulesDir) {
			return filepath.SkipDir
		}

		// Only process .yaml and .yml files
		if info.IsDir() || (!strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml")) {
			return nil
		}

		// #nosec G304 - path is validated to be within rulesDir
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var rule Rule
		if err := yaml.Unmarshal(data, &rule); err != nil {
			return err
		}

		rules = append(rules, rule)
		return nil
	})

	return rules, err
}
