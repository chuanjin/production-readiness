package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/output"
	"github.com/chuanjin/production-readiness/internal/rules"
	"github.com/chuanjin/production-readiness/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	format string
	debug  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a codebase and evaluate production readiness",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println("Invalid path:", err)
			return
		}

		// 1️⃣ load rules
		ruleSet, err := rules.LoadRules("rules")
		if err != nil {
			fmt.Println("Error loading rules:", err)
			return
		}

		// 2️⃣ scan repo with debug option
		var logger scanner.Logger = &scanner.NoopLogger{}
		if debug {
			logger = log.New(os.Stdout, "[scanner] ", log.LstdFlags)
		}

		signals, err := scanner.ScanRepoWithOptions(absPath, scanner.ScanOptions{
			Debug:  debug,
			Logger: logger,
		})
		if err != nil {
			fmt.Println("Error scanning:", err)
			return
		}

		// 3️⃣ evaluate
		findings := engine.Evaluate(ruleSet, signals)

		// Summarize
		summary := engine.Summarize(findings)

		// 4️⃣ output
		switch format {
		case "json":
			fmt.Println(output.JSON(summary, findings, &signals))
		default:
			fmt.Println(output.Markdown(summary, findings, &signals))
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&format, "format", "f", "md", "output format: md or json")
	scanCmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
}
