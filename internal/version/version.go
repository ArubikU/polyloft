package version

import "runtime/debug"

// These variables can be overridden at build time via -ldflags.
var (
	// Version is the semantic version or commit hash.
	Version = "v1.1.6"
	// Commit is the VCS commit identifier.
	Commit = ""
	// Date is the build date.
	Date = "2025-10-29"
)

// String returns a human-friendly version string.
func String() string {
	v := Version
	if info, ok := debug.ReadBuildInfo(); ok && v == "" {
		v = info.Main.Version
	}
	if Commit != "" {
		v += " (" + Commit[:min(7, len(Commit))] + ")"
	}
	if Date != "" {
		v += " built " + Date
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
