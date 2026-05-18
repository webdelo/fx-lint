package fxlint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	fxlint "github.com/webdelo/fx-lint"
)

func TestModuleLocationAnalyzer(t *testing.T) {
	t.Parallel()
	a := fxlint.NewModuleLocationAnalyzer(fxlint.ModuleLocationConfig{ModPath: "example.com/proj"})
	analysistest.Run(t, analysistest.TestData(), a,
		"bad/outside_modules",
	)
}

func TestModuleLocationAnalyzer_FxImportOutsideTree(t *testing.T) {
	t.Parallel()
	a := fxlint.NewModuleLocationAnalyzer(fxlint.ModuleLocationConfig{ModPath: "example.com/proj"})
	analysistest.Run(t, analysistest.TestData(), a,
		"bad/outside_fx_import",
	)
}

func TestModuleLocationAnalyzer_Clean(t *testing.T) {
	t.Parallel()
	a := fxlint.NewModuleLocationAnalyzer(fxlint.ModuleLocationConfig{ModPath: "example.com/proj"})
	analysistest.Run(t, analysistest.TestData(), a,
		"example.com/proj/internal/infrastructure/fx/modules/fxm_example",
	)
}

func TestModuleLocationAnalyzer_AppHelpersBoundary(t *testing.T) {
	t.Parallel()
	a := fxlint.NewModuleLocationAnalyzer(fxlint.ModuleLocationConfig{ModPath: "example.com/proj"})
	analysistest.Run(t, analysistest.TestData(), a,
		"example.com/proj/internal/app_helpers",
	)
}

func TestModuleLocationAnalyzer_EmptyModPathPanics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty ModPath")
		}
	}()
	fxlint.NewModuleLocationAnalyzer(fxlint.ModuleLocationConfig{})
}
