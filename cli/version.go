package cli

import (
	"fmt"

	"github.com/chuanjin/production-readiness/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Production-Readiness CLI %s (%s) built %s\n",
			version.Version, version.Commit, version.BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
