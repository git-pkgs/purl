package purl

import (
	"github.com/git-pkgs/vers"
)

// CleanVersion extracts a version from a version constraint string.
// Uses the vers library to parse the constraint and extract the minimum bound.
// If parsing fails, returns the original string.
func CleanVersion(version, scheme string) string {
	if version == "" {
		return ""
	}

	r, err := vers.ParseNative(version, scheme)
	if err != nil || len(r.Intervals) == 0 {
		return version
	}

	// Return the minimum bound from the first interval
	if r.Intervals[0].Min != "" {
		return r.Intervals[0].Min
	}

	return version
}
