// Package fx is a minimal stub of go.uber.org/fx for use in analysistest
// testdata packages.  Only the symbols exercised by fxlint tests are declared.
package fx

import "context"

// Option is a functional option for fx.New.
type Option interface{}

// Provide registers constructors with the DI container.
func Provide(constructors ...interface{}) Option { return nil }

// Invoke registers functions that are run after the container is built.
func Invoke(funcs ...interface{}) Option { return nil }

// Module groups a set of options under a named scope.
func Module(name string, opts ...Option) Option { return nil }

// Options groups a set of options without module semantics.
func Options(opts ...Option) Option { return nil }

// Hook is a pair of OnStart/OnStop callbacks.
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// Lifecycle accepts hooks that are invoked at app start/stop.
type Lifecycle interface {
	Append(Hook)
}
