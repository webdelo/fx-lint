// Package fxm_constructors — providers.go illegally defines a constructor.
package fxm_constructors

// newSomething is a constructor that must NOT appear in providers.go of an
// fxm_* package (UFX-002/UFX-024).
func newSomething() *struct{} { return nil } // want `constructor "newSomething" in fxm_\*/providers\.go violates UFX-002`
