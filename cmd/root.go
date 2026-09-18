package cmd

import (
	"fmt"
	"os"

	"github.com/padok-team/guacamole-pre-commit/internal/pipeline"
	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd is the entrypoint invoked by the pre-commit framework, which calls
// the binary with the staged files (.tf/.hcl) as positional arguments
var rootCmd = &cobra.Command{
	Use:   "guacamole-pre-commit [files...]",
	Short: "Run Guacamole's static checks on staged Terraform/Terragrunt files.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		exitCode, err := pipeline.Run(args, os.Stdout, verbose)
		if err != nil {
			fmt.Fprintln(os.Stderr, "guacamole-pre-commit:", err)
			os.Exit(1)
		}
		os.Exit(exitCode)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print every file received, every module/layer dir scoped, and every check run (pass or fail) — not just failures")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
