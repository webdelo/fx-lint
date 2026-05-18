package fxlint

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ModuleLocationConfig holds configuration for the ModuleLocation analyzer.
// ModPath must be set to the Go module path of the project being linted
// (e.g. "github.com/myorg/myrepo"). It is used to derive the allowed import
// and call prefixes so that the analyzer is not hard-coded to any single
// repository.
type ModuleLocationConfig struct {
	ModPath string
}

// NewModuleLocationAnalyzer returns a configured ModuleLocation analyzer.
// It reports two kinds of wiring-leak violations:
//
//  1. Direct calls to fx.Module / fx.Provide / fx.Invoke from a package outside
//     the designated wiring tree.
//  2. Any import of go.uber.org/fx from a package outside the designated wiring
//     tree.  This catches wiring code that does not call fx.* directly but still
//     refers to fx.Lifecycle / fx.Hook / fx.In / fx.Out / fx.Annotate (typical
//     for misplaced lifecycle hooks and provider wrappers).
//
// Only packages rooted under the allowed prefixes may import or call into
// go.uber.org/fx; all other locations indicate wiring logic that leaked out of
// the infrastructure layer.
//
// Panics if cfg.ModPath is empty — this is a programmer-error guard; configuration
// validation belongs in callers (e.g. the plugin shim or cmd/fxlint) that surface
// errors to end users before reaching this constructor.
func NewModuleLocationAnalyzer(cfg ModuleLocationConfig) *analysis.Analyzer {
	if cfg.ModPath == "" {
		panic("fxlint: ModuleLocationConfig.ModPath must be non-empty")
	}

	// allowedImportPrefixes lists package path prefixes where importing
	// go.uber.org/fx is permitted.  This covers the entire infrastructure/fx tree
	// because shows (fxs_*) may need fx.In/fx.Out and hooks (fxh_*) need
	// fx.Lifecycle/fx.Hook.  app/ and cmd/ are the entry-point roots.
	allowedImportPrefixes := []string{
		cfg.ModPath + "/internal/infrastructure/fx",
		cfg.ModPath + "/internal/app",
		cfg.ModPath + "/internal/cmd",
	}

	// allowedCallPrefixes lists package path prefixes where direct calls to
	// fx.Module / fx.Provide / fx.Invoke are permitted.  These are wiring-only
	// decisions and must not appear in shows (fxs_*), hooks (fxh_*), adapters
	// (fxa_*), or any non-fx package.
	allowedCallPrefixes := []string{
		cfg.ModPath + "/internal/infrastructure/fx/modules",
		cfg.ModPath + "/internal/app",
		cfg.ModPath + "/internal/cmd",
	}

	return &analysis.Analyzer{
		Name:     "fxmodulelocation",
		Doc:      "reports fx.Module/Provide/Invoke calls and any go.uber.org/fx import outside the allowed wiring packages",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return runModuleLocation(pass, allowedImportPrefixes, allowedCallPrefixes)
		},
	}
}

func runModuleLocation(pass *analysis.Pass, allowedImportPrefixes, allowedCallPrefixes []string) (interface{}, error) {
	pkgPath := pass.Pkg.Path()

	if !hasPrefix(pkgPath, allowedImportPrefixes) {
		reportFxImports(pass)
	}

	if !hasPrefix(pkgPath, allowedCallPrefixes) {
		reportFxCalls(pass)
	}

	return nil, nil
}

// reportFxImports flags any non-test file that imports go.uber.org/fx from a
// package outside the wiring tree.  The diagnostic is attached to the import
// spec so that the offending file is unambiguous.
func reportFxImports(pass *analysis.Pass) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if isTestFile(filename) {
			continue
		}

		for _, imp := range file.Imports {
			if imp.Path == nil {
				continue
			}

			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}

			if path != fxPkgPath {
				continue
			}

			pass.Reportf(imp.Pos(),
				"import of %q from %q is outside the allowed wiring tree (internal/infrastructure/fx/, internal/app/, internal/cmd/); fx wiring (Module/Provide/Invoke, Lifecycle/Hook, In/Out/Annotate) must live in fxm_*/fxs_*/fxh_*/fxa_* packages",
				fxPkgPath, pass.Pkg.Path())
		}
	}
}

// reportFxCalls flags direct calls to fx.Module/fx.Provide/fx.Invoke.  This is
// kept as a separate, more specific diagnostic because it pinpoints the exact
// call site even in packages that legitimately import fx (for instance via a
// re-export elsewhere in the future).
func reportFxCalls(pass *analysis.Pass) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Skip test files: fx.ValidateApp stubs in _test.go files are not
		// production wiring and do not violate UFX-024.
		if pos := pass.Fset.Position(call.Pos()); isTestFile(pos.Filename) {
			return
		}

		name, ok := fxCallName(pass, call)
		if !ok {
			return
		}

		if name != "Module" && name != "Provide" && name != "Invoke" {
			return
		}

		pass.Reportf(call.Pos(),
			"fx.%s called from %q which is outside the allowed wiring packages; move to internal/infrastructure/fx/modules/",
			name, pass.Pkg.Path())
	})
}

// hasPrefix reports whether pkgPath equals one of the prefixes or starts with
// "<prefix>/".  This avoids accidental matches such as ".../app" matching a
// hypothetical ".../app_helpers" package.
func hasPrefix(pkgPath string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if pkgPath == prefix || strings.HasPrefix(pkgPath, prefix+"/") {
			return true
		}
	}

	return false
}
