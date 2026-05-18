// Package fxm_inline_lambda contains fx.Provide/Invoke calls with inline lambdas.
package fxm_inline_lambda

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("bad",
		fx.Provide(func() *struct{} { return nil }), // want `inline lambda in fx\.Provide violates UFX-026`
		fx.Invoke(func(*struct{}) {}),               // want `inline lambda in fx\.Invoke violates UFX-026`
	)
}
