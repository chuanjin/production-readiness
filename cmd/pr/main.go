// cmd/pr/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/output"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		fmt.Println("Usage: pr scan [path]")
		os.Exit(1)
	}

	root := "."
	if len(os.Args) >= 3 {
		root = os.Args[2]
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Println("Invalid path:", err)
		os.Exit(1)
	}

	// Load rules
	ruleSet, err := rules.LoadFromDir("rules")
	if err != nil {
		fmt.Println("Failed to load rules:", err)
		os.Exit(1)
	}

	// Scan repository
	signals, err := scanner.ScanRepo(absRoot)
	if err != nil {
		fmt.Println("Failed to scan repo:", err)
		os.Exit(1)
	}

	// Evaluate rules  => returns engine.Summary (struct)
	summary := engine.Evaluate(ruleSet, signals)

	// Output Markdown (only one argument)
	report := output.Markdown(summary)

	fmt.Println(report)
}
