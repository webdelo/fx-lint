package fxlint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const fxPkgPath = "go.uber.org/fx"

// fxCallName returns the selector name (e.g. "Provide") and true when call is
// a direct call to a function in the go.uber.org/fx package.  It uses type
// information to resolve the package so that a local variable named "fx" cannot
// produce false positives.
func fxCallName(pass *analysis.Pass, call *ast.CallExpr) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	obj := objectOf(pass, sel)
	if obj == nil {
		return "", false
	}

	pkg := obj.Pkg()
	if pkg == nil || pkg.Path() != fxPkgPath {
		return "", false
	}

	return sel.Sel.Name, true
}

// objectOf resolves the object referred to by the selector expression using
// type information.  It handles both qualified identifiers (pkg.Func) and
// method expressions, preferring Uses over Selections so that package-level
// functions are found correctly.
func objectOf(pass *analysis.Pass, sel *ast.SelectorExpr) types.Object {
	// Package-qualified identifier: the selector Sel is looked up via Uses.
	if obj, ok := pass.TypesInfo.Uses[sel.Sel]; ok {
		return obj
	}

	// Method / field selection on a value.
	if selection, ok := pass.TypesInfo.Selections[sel]; ok {
		return selection.Obj()
	}

	return nil
}
