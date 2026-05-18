package fxlint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// InlineLambdaAnalyzer reports inline function literals passed directly to
// fx.Provide or fx.Invoke (UFX-026 violation).  Every constructor must be a
// named function so that fx's reflection-based debug output and the dependency
// graph are legible.
var InlineLambdaAnalyzer = &analysis.Analyzer{
	Name:     "fxinlinelambda",
	Doc:      "reports inline lambdas in fx.Provide/fx.Invoke (UFX-026)",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runInlineLambda,
}

func runInlineLambda(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Skip positions in test files (_test.go): UFX-026 targets production wiring
		// code only; test stubs legitimately use inline provider functions.
		if pos := pass.Fset.Position(call.Pos()); isTestFile(pos.Filename) {
			return
		}

		name, ok := fxCallName(pass, call)
		if !ok {
			return
		}

		if name != "Provide" && name != "Invoke" {
			return
		}

		for _, arg := range call.Args {
			if _, isLit := arg.(*ast.FuncLit); isLit {
				pass.Reportf(arg.Pos(),
					"inline lambda in fx.%s violates UFX-026; extract to a named function in fxs_*/fxh_*",
					name)
			}
		}
	})

	return nil, nil
}

// isTestFile reports whether the given source file path is a Go test file.
func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}
