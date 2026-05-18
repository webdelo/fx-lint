// Package fxm_clean is an fxm_* module.go that correctly contains only the
// Module orchestration function with a reference to a named constructor.
// The named constructor lives in a separate provider file (not shown here).
package fxm_clean

import "go.uber.org/fx"

// Module wires the clean sub-system.  Only the Module function is defined here.
func Module() fx.Option {
	return fx.Module("clean",
		fx.Provide(newCleanService),
	)
}
