# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-18

### Added
- Initial release.
- `InlineLambdaAnalyzer` (UFX-026): reports inline lambdas in `fx.Provide`/`fx.Invoke`.
- `ModuleLocationAnalyzer` (UFX-024): reports `fx.Module`/`Provide`/`Invoke` calls and `go.uber.org/fx` imports outside the allowed wiring tree. Configurable via `ModuleLocationConfig{ModPath string}`.
- `ModuleOrchestrationAnalyzer` (UFX-002): reports constructor definitions inside `fxm_*/module.go|providers.go|lifecycle.go`.
- `plugin/` subpackage implementing `register.LinterPlugin` for golangci-lint v2 custom builds.
- `cmd/fxlint` standalone multichecker with `-mod-path` flag.
