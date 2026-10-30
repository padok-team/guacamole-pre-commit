package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/gruntwork-io/terragrunt/config"
	"github.com/gruntwork-io/terragrunt/options"
)

// Scope climbs from each received file's directory to the nearest enclosing
// Terraform module dir (contains *.tf) or Terragrunt layer dir (contains
// terragrunt.hcl). For a layer, it also resolves the module it calls via
// module.hcl's `terraform { source = ... }` block, if present. A path whose
// climb reaches the filesystem root without finding either marker is out of
// IaC scope and is skipped silently. Returns deduplicated, sorted, absolute
// directories.
func Scope(paths []string) (moduleDirs, layerDirs []string, err error) {
	modules := map[string]struct{}{}
	layers := map[string]struct{}{}

	for _, p := range paths {
		abs, absErr := filepath.Abs(p)
		if absErr != nil {
			return nil, nil, absErr
		}

		dir := filepath.Dir(abs)
		for {
			if isLayerDir(dir) {
				layers[dir] = struct{}{}
				if moduleDir, ok := resolveModuleSource(dir); ok {
					modules[moduleDir] = struct{}{}
				}
				break
			}
			if isModuleDir(dir) {
				modules[dir] = struct{}{}
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break // filesystem root reached, out of IaC scope
			}
			dir = parent
		}
	}

	return sortedKeys(modules), sortedKeys(layers), nil
}

func isModuleDir(dir string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, "*.tf"))
	return err == nil && len(matches) > 0
}

func isLayerDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "terragrunt.hcl"))
	return err == nil
}

// resolveModuleSource finds the module.hcl governing layerDir and reads its
// `terraform { source = ... }` block, evaluated with Terragrunt's own parser
// (so functions like get_repo_root()/get_path_to_repo_root() resolve), and
// returns the module dir it points to. Best effort: no module.hcl found, a
// parse error, or a source that isn't a local path with .tf files just mean
// no module is added for this layer, not an error.
func resolveModuleSource(layerDir string) (string, bool) {
	moduleHCL, ok := findModuleHCL(layerDir)
	if !ok {
		return "", false
	}

	opts, err := options.NewTerragruntOptionsWithConfigPath(moduleHCL)
	if err != nil {
		return "", false
	}

	parsingCtx := config.NewParsingContext(context.Background(), opts).WithDecodeList(config.TerraformSource)
	cfg, err := config.PartialParseConfigFile(parsingCtx, moduleHCL, nil)
	if err != nil || cfg.Terraform == nil || cfg.Terraform.Source == nil || *cfg.Terraform.Source == "" {
		return "", false
	}

	source := *cfg.Terraform.Source
	if !filepath.IsAbs(source) {
		// Relative sources (including what get_path_to_repo_root() returns)
		// are relative to the directory that actually contains module.hcl,
		// not necessarily the leaf layer that included it.
		source = filepath.Join(filepath.Dir(moduleHCL), source)
	}
	source = filepath.Clean(source)

	if !isModuleDir(source) {
		return "", false
	}
	return source, true
}

// findModuleHCL mirrors Terragrunt's own find_in_parent_folders("module.hcl"):
// module.hcl is often shared across several leaf layers higher up the tree
// (like root.hcl), not duplicated in every layer directory, so this climbs
// from dir through its ancestors until a module.hcl is found or the
// filesystem root is reached.
func findModuleHCL(dir string) (string, bool) {
	for {
		candidate := filepath.Join(dir, "module.hcl")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
