package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the entrypoint invoked by the pre-commit framework, which calls
// the binary with the staged files (.tf/.hcl) as positional arguments
var rootCmd = &cobra.Command{
	Use:   "guacamole-pre-commit [files...]",
	Short: "Run Guacamole's static checks on staged Terraform/Terragrunt files.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "guacamole-pre-commit: no files received")
			return
		}
		fmt.Printf("guacamole-pre-commit: received %d file(s):\n", len(args))
		for _, f := range args {
			fmt.Println("  -", f)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
