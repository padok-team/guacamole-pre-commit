package pipeline

import (
	"fmt"
	"io"
	"sort"

	"github.com/padok-team/guacamole/helpers"
)

// Run scopes paths to their module/layer dirs, runs the applicable
// Guacamole static checks, and writes any errors to w as
// "path:line [ID] message" (LineNumber == -1 omits ":line"). Returns 0 if
// every check passed, 1 otherwise, and writes nothing when everything
// passes.
//
// When verbose is true, it additionally logs the files received and the
// module/layer dirs they were scoped to on w, then renders every check that
// ran (pass or fail) via Guacamole's own helpers.RenderChecks — the same
// "✅ ID - Name" / "  -  path:line" / "Score: X% (Y/Z)" format the
// guacamole CLI itself uses, so it's the format users already know. That
// helper writes straight to the process's real stdout (not w — it's not
// parameterized upstream), which only matters if this is ever called with a
// non-os.Stdout writer.
func Run(paths []string, w io.Writer, verbose bool) (int, error) {
	if verbose {
		fmt.Fprintf(w, "guacamole-pre-commit: %d file(s) received:\n", len(paths))
		for _, p := range paths {
			fmt.Fprintln(w, "  -", p)
		}
	}

	moduleDirs, layerDirs, err := Scope(paths)
	if err != nil {
		return 1, err
	}

	if verbose {
		fmt.Fprintf(w, "guacamole-pre-commit: scoped to %d module dir(s), %d layer dir(s):\n", len(moduleDirs), len(layerDirs))
		for _, d := range moduleDirs {
			fmt.Fprintln(w, "  module:", d)
		}
		for _, d := range layerDirs {
			fmt.Fprintln(w, "  layer: ", d)
		}
	}

	allChecks := append(runModuleChecks(moduleDirs), runLayerChecks(layerDirs)...)
	sort.Slice(allChecks, func(i, j int) bool { return allChecks[i].ID < allChecks[j].ID })

	exitCode := 0
	for _, check := range allChecks {
		if len(check.Errors) > 0 {
			exitCode = 1
		}
	}

	if verbose {
		if len(allChecks) > 0 {
			helpers.RenderChecks(allChecks, true)
		}
		return exitCode, nil
	}

	for _, check := range allChecks {
		for _, e := range check.Errors {
			if e.LineNumber == -1 {
				fmt.Fprintf(w, "%s [%s] %s\n", e.Path, check.ID, e.Description)
			} else {
				fmt.Fprintf(w, "%s:%d [%s] %s\n", e.Path, e.LineNumber, check.ID, e.Description)
			}
		}
	}
	return exitCode, nil
}
