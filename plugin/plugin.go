// Package plugin registers fxlint as a golangci-lint v2 custom plugin.
//
// Usage in .golangci.yaml:
//
//	linters:
//	  settings:
//	    custom:
//	      fxlint:
//	        type: module
//	        settings:
//	          mod_path: github.com/your/project
package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	fxlint "github.com/webdelo/fx-lint"
)

func init() { register.Plugin("fxlint", New) }

// Settings is the JSON-deserializable plugin configuration block.
type Settings struct {
	ModPath string `json:"mod_path"`
}

type plugin struct {
	settings Settings
}

// New constructs the fxlint plugin instance from raw golangci-lint settings.
// The settings block is converted via a JSON round-trip because golangci-lint
// passes settings as `any` decoded from YAML into map[string]interface{}.
func New(raw any) (register.LinterPlugin, error) {
	var s Settings
	if raw != nil {
		b, err := json.Marshal(raw)
		if err != nil {
			return nil, fmt.Errorf("fxlint: marshal settings: %w", err)
		}
		if err := json.Unmarshal(b, &s); err != nil {
			return nil, fmt.Errorf("fxlint: unmarshal settings: %w", err)
		}
	}
	if s.ModPath == "" {
		return nil, fmt.Errorf("fxlint: settings.mod_path is required")
	}
	return &plugin{settings: s}, nil
}

// BuildAnalyzers returns the configured fxlint analyzers.
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return fxlint.NewAnalyzers(fxlint.ModuleLocationConfig{ModPath: p.settings.ModPath}), nil
}

// GetLoadMode reports the package load mode required by fxlint analyzers.
// LoadModeTypesInfo is required because fxCallName resolves the fx selector
// via pass.TypesInfo.Uses.
func (p *plugin) GetLoadMode() string { return register.LoadModeTypesInfo }
