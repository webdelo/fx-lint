package fxlint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	fxlint "github.com/webdelo/fx-lint"
)

func TestInlineLambdaAnalyzer(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), fxlint.InlineLambdaAnalyzer,
		"bad/fxm_inline_lambda",
	)
}

func TestInlineLambdaAnalyzer_Clean(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), fxlint.InlineLambdaAnalyzer,
		"good/fxm_clean",
		"good/named_provider",
	)
}
