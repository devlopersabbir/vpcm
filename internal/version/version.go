// Package version holds the single source of truth for the VPSM version.
// The Version variable is injected at build time via -ldflags:
//
//	go build -ldflags "-X github.com/devlopersabbir/vpcm/internal/version.Version=v1.2.0"
//
// In local development it falls back to the value set here.
package version

// Version is bumped automatically by the CI release workflow.
// Only the patch number increments per release (e.g. v1.1.2 → v1.1.3).
// Do NOT edit manually — CI reads this value, increments the patch, and
// commits the updated file back to main after every release.
var Version = "v1.1.13"
