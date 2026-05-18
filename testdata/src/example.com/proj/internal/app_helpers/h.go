// Package app_helpers must NOT be matched by the allowed prefix
// "example.com/proj/internal/app" — the hasPrefix guard ensures only exact
// equality or "<prefix>/..." matches.  Importing fx from here is a violation.
package app_helpers

import "go.uber.org/fx" // want `import of "go.uber.org/fx" from "example.com/proj/internal/app_helpers" is outside the allowed wiring tree`

var _ = fx.Provide
