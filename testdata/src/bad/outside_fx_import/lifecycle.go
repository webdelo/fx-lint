// Package outside_fx_import is a non-wiring package that imports go.uber.org/fx
// to define a lifecycle hook.  The package does not call fx.Module/Provide/Invoke
// directly (so the call-level rule alone would miss it), but importing fx from
// outside the wiring tree is itself a structural violation: lifecycle hooks
// belong in fxh_*, provider wrappers in fxs_*.
package outside_fx_import

import (
	"context"

	"go.uber.org/fx" // want `import of "go.uber.org/fx" from "bad/outside_fx_import" is outside the allowed wiring tree`
)

// RegisterLifecycle is shaped like a typical fx hook (lc fx.Lifecycle, deps...)
// but lives outside the fxh_* tree.  This is exactly the misplacement pattern
// observed in internal/infrastructure/outbox/guard/lifecycle.go.
func RegisterLifecycle(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error { return nil },
		OnStop:  func(_ context.Context) error { return nil },
	})
}
