// Command fxlint is a standalone go/analysis multichecker that runs the
// fxlint analyzers against the supplied package patterns.
//
// Usage:
//
//	fxlint -mod-path=github.com/your/project ./...
//
// The -mod-path flag (or FX_LINT_MOD_PATH environment variable) is required
// and specifies the module path of the target project — it determines the
// whitelist prefixes used by the ModuleLocation analyzer.
//
// All other flags (-flags, -json, -c, -fix, -diff, -V, plus per-analyzer
// flags) are forwarded to multichecker.Main unchanged.
package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis/multichecker"

	fxlint "github.com/webdelo/fx-lint"
)

func main() {
	modPath, rest := extractModPath(os.Args[1:], os.Getenv("FX_LINT_MOD_PATH"))
	if modPath == "" {
		log.Fatal("fxlint: -mod-path (or FX_LINT_MOD_PATH) is required")
	}
	os.Args = append([]string{os.Args[0]}, rest...)
	multichecker.Main(fxlint.NewAnalyzers(fxlint.ModuleLocationConfig{ModPath: modPath})...)
}

// extractModPath consumes -mod-path / --mod-path (with or without =value) from
// args and returns its value (or fallback) plus the remaining args untouched.
// Other flags pass through to multichecker.Main which handles them itself.
func extractModPath(args []string, fallback string) (string, []string) {
	value := fallback
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-mod-path", a == "--mod-path":
			if i+1 < len(args) {
				value = args[i+1]
				i++
			}
		case strings.HasPrefix(a, "-mod-path="):
			value = strings.TrimPrefix(a, "-mod-path=")
		case strings.HasPrefix(a, "--mod-path="):
			value = strings.TrimPrefix(a, "--mod-path=")
		default:
			out = append(out, a)
		}
	}
	return value, out
}
