// Package fxm_clean — lifecycle.go correctly contains only the Lifecycle
// orchestration function. No constructors belong here.
package fxm_clean

import "go.uber.org/fx"

// Lifecycle wires fx lifecycle hooks for the clean sub-system.
func Lifecycle() fx.Option {
	return fx.Options()
}
