package pipeline

import (
	"os"

	"github.com/padok-team/guacamole/checks"
	"github.com/padok-team/guacamole/data"
	"github.com/spf13/viper"
)

// runScoped calls runAll once per dir, pointing Guacamole's global
// "codebase-path" viper key at that single directory each time so runAll
// (checks.ModuleStaticChecks or checks.LayerStaticChecks) only ever sees
// that dir's files, then merges the returned checks by ID across
// iterations. This runs every check Guacamole ships, with no per-check
// list to keep in sync as guacamole gains checks.
func runScoped(dirs []string, runAll func() []data.Check) []data.Check {
	previous := viper.GetString("codebase-path")
	defer viper.Set("codebase-path", previous)

	merged := map[string]*data.Check{}
	order := []string{}

	for _, dir := range dirs {
		viper.Set("codebase-path", dir)
		for _, check := range runAll() {
			if existing, ok := merged[check.ID]; ok {
				existing.Errors = append(existing.Errors, check.Errors...)
				if len(existing.Errors) > 0 {
					existing.Status = "❌"
				}
				continue
			}
			c := check
			merged[check.ID] = &c
			order = append(order, check.ID)
		}
	}

	result := make([]data.Check, 0, len(order))
	for _, id := range order {
		result = append(result, *merged[id])
	}
	return result
}

func runModuleChecks(moduleDirs []string) []data.Check {
	return runScoped(moduleDirs, checks.ModuleStaticChecks)
}

// runLayerChecks runs the layer checks away from the repo's working
// directory. Guacamole's Dry check derives its default Terragrunt download
// dir from the process CWD and skips any scanned path found nested under
// it (see tests/README.md's CWD trap) — since pre-commit always invokes
// this binary with CWD at the repo root, a parent of every layer dir, that
// default silently skips every layer unless CWD is moved elsewhere first.
func runLayerChecks(layerDirs []string) []data.Check {
	if len(layerDirs) == 0 {
		return nil
	}

	cwd, err := os.Getwd()
	if err == nil {
		if err := os.Chdir(os.TempDir()); err == nil {
			defer os.Chdir(cwd)
		}
	}

	return runScoped(layerDirs, checks.LayerStaticChecks)
}
