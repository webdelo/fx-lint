// Package fxm_example is a correctly placed fx module inside the allowed tree.
package fxm_example

import "go.uber.org/fx"

// Module wires the example sub-system.
func Module() fx.Option {
	return fx.Module("example", fx.Provide(newExampleService))
}
