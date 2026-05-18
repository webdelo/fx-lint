// Package fxm_constructors is an fxm_* module.go that illegally defines constructors.
package fxm_constructors

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("bad", fx.Provide(newService))
}

func newService() *struct{} { return nil } // want `constructor "newService" in fxm_\*/module\.go violates UFX-002`
