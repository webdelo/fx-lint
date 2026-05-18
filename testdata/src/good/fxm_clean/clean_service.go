// Package fxm_clean — clean_service.go holds the constructor for the clean
// sub-system. This file is NOT named module.go/providers.go/lifecycle.go so it
// is exempt from the orchestration-only lint rule (UFX-002/UFX-024).
package fxm_clean

// newCleanService is a named constructor living in a non-orchestration file.
// This is the correct location per UFX-002.
func newCleanService() *struct{} { return nil }
