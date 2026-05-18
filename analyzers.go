// Package fxlint provides go/analysis linters that enforce project-specific
// conventions for uber-go/fx usage (UFX-002/024/026).
//
// Two analyzers are stateless and exported as package variables
// (InlineLambdaAnalyzer, ModuleOrchestrationAnalyzer). The third
// (ModuleLocationAnalyzer) requires a runtime ModPath and is built via
// NewModuleLocationAnalyzer. Use NewAnalyzers to obtain the full set.
package fxlint

import "golang.org/x/tools/go/analysis"

// NewAnalyzers returns the full set of fxlint analyzers configured with cfg.
// The returned slice is safe to pass to multichecker.Main or to a
// golangci-lint plugin BuildAnalyzers result.
func NewAnalyzers(cfg ModuleLocationConfig) []*analysis.Analyzer {
	return []*analysis.Analyzer{
		InlineLambdaAnalyzer,
		NewModuleLocationAnalyzer(cfg),
		ModuleOrchestrationAnalyzer,
	}
}
