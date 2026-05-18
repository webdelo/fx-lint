package fxlint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	fxlint "github.com/webdelo/fx-lint"
)

func TestModuleOrchestrationAnalyzer(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), fxlint.ModuleOrchestrationAnalyzer,
		"bad/fxm_constructors",
	)
}

func TestModuleOrchestrationAnalyzer_Clean(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), fxlint.ModuleOrchestrationAnalyzer,
		"good/fxm_clean",
	)
}
