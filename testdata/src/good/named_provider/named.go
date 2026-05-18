// Package named_provider demonstrates fx.Provide with a named constructor — no violation.
package named_provider

import "go.uber.org/fx"

func newService() *struct{} { return nil }

// Providers returns the fx options for this package.
func Providers() fx.Option {
	return fx.Provide(newService)
}
