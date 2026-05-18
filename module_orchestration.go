package fxlint

import (
	"go/ast"
	"path"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ModuleOrchestrationAnalyzer enforces that the orchestration files inside
// fxm_* packages (module.go, providers.go, lifecycle.go) contain only the
// single named entry-point function each file is responsible for.
// Constructor functions (newXxx, provideXxx, registerXxx, adapt*, and any
// method receiver declarations) are violations of UFX-002/UFX-024 because
// they introduce construction logic that belongs in fxs_*/fxh_* provider files.
var ModuleOrchestrationAnalyzer = &analysis.Analyzer{
	Name:     "fxmoduleorchestration",
	Doc:      "reports constructor definitions inside fxm_*/module.go|providers.go|lifecycle.go (UFX-002/UFX-024)",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runModuleOrchestration,
}

// allowedModuleFiles is the set of file names inside fxm_* packages that are
// subject to the orchestration-only lint rule.
var allowedModuleFiles = map[string]bool{
	"module.go":    true,
	"providers.go": true,
	"lifecycle.go": true,
}

// allowedTopLevelFuncs are the only function names permitted at the top level
// of an orchestration file inside an fxm_* package.
var allowedTopLevelFuncs = map[string]bool{
	"Module":    true,
	"Providers": true,
	"Lifecycle": true,
}

func runModuleOrchestration(pass *analysis.Pass) (interface{}, error) {
	pkgPath := pass.Pkg.Path()
	// Only enforce inside fxm_* packages.
	pkgName := path.Base(pkgPath)
	if !strings.HasPrefix(pkgName, "fxm_") {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.File)(nil)}, func(n ast.Node) {
		file := n.(*ast.File)
		baseName := path.Base(pass.Fset.File(file.Pos()).Name())
		if !allowedModuleFiles[baseName] {
			return
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Method receivers do not belong in orchestration files.
			if fn.Recv != nil {
				pass.Reportf(fn.Pos(),
					"method declaration in fxm_*/%s violates UFX-002; move to a provider file",
					baseName)
				continue
			}

			if !allowedTopLevelFuncs[fn.Name.Name] {
				pass.Reportf(fn.Pos(),
					"constructor %q in fxm_*/%s violates UFX-002/UFX-024; orchestration files must be orchestration only — move to fxs_*/fxh_*",
					fn.Name.Name, baseName)
			}
		}
	})

	return nil, nil
}
