// Package outside_modules calls fx.Module from a package path that is not
// inside the allowed wiring tree — a location violation.
package outside_modules

import "go.uber.org/fx" // want `import of "go.uber.org/fx" from "bad/outside_modules" is outside the allowed wiring tree`

func Register() fx.Option {
	opts := fx.Provide(newThing)      // want `fx\.Provide called from`
	return fx.Module("outside", opts) // want `fx\.Module called from`
}

func newThing() *struct{} { return nil }
