package rules

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadFromDir(dir string) ([]Rule, error) {
	var rules []Rule

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}

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
