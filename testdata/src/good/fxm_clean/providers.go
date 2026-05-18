// Package fxm_clean — providers.go correctly contains only the Providers
// orchestration function. No constructors belong here.
package fxm_clean

import "go.uber.org/fx"

// Providers returns the fx.Option set for the clean sub-system.
func Providers() fx.Option {
	return fx.Provide(newCleanService)
}
