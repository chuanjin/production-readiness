package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pr",
	Short: "Production-Readiness CLI — audit repos for deployability & resilience",
	Long: `Production-Readiness is a senior-engineering-informed static scanner 
that identifies operational, resiliency, observability and rollback risks 
before you deploy to production.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
