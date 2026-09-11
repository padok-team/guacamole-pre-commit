package cmd

import (
	"fmt"
	"os"

	"github.com/padok-team/guacamole/checks"
	"github.com/padok-team/guacamole/data"
	"github.com/padok-team/guacamole/helpers"
	"github.com/spf13/cobra"
)

// pocCmd: prove guacamole-pre-commit can import Guacamole as
// a library and run one of its checks, without reimplementing the check logic.
// Throwaway once T2.2/T2.3 wire real file->module scoping and check execution.
var pocCmd = &cobra.Command{
	Use:    "poc [module-dir]",
	Short:  "[POC T2.1] Run Guacamole's VarContainsDescription check on a single module directory",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		module, err := helpers.LoadModule(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "load module:", err)
			os.Exit(1)
		}

		modules := map[string]data.TerraformModule{path: module}
		check, err := checks.VarContainsDescription(modules)
		if err != nil {
			fmt.Fprintln(os.Stderr, "run check:", err)
			os.Exit(1)
		}

		fmt.Printf("%s %s\n", check.Status, check.ID)
		for _, e := range check.Errors {
			fmt.Printf("  - %s:%d %s\n", e.Path, e.LineNumber, e.Description)
		}

		if check.Status == "❌" {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pocCmd)
}
