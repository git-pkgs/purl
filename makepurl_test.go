package purl

import (
	"testing"
)

func TestCleanVersion(t *testing.T) {
	tests := []struct {
		version string
		scheme  string
		want    string
	}{
		// npm constraints
		{"1.0.0", "npm", "1.0.0"},
		{"^1.0.0", "npm", "1.0.0"},
		{"~1.0.0", "npm", "1.0.0"},
		{">=1.0.0", "npm", "1.0.0"},
		{">=1.0.0 <2.0.0", "npm", "1.0.0"},

		// gem constraints
		{"~> 1.0", "gem", "1.0"},
		{">= 1.0, < 2.0", "gem", "1.0"},

		// pypi constraints
		{">=1.0.0", "pypi", "1.0.0"},
		{"~=1.4.2", "pypi", "1.4.2"},

		// cargo constraints
		{"^1.0.0", "cargo", "1.0.0"},

		// Plain versions pass through
		{"1.0.0", "npm", "1.0.0"},
		{"v1.0.0", "go", "v1.0.0"},

		// Empty
		{"", "npm", ""},
	}

	for _, tt := range tests {
		t.Run(tt.version+"_"+tt.scheme, func(t *testing.T) {
			if got := CleanVersion(tt.version, tt.scheme); got != tt.want {
				t.Errorf("CleanVersion(%q, %q) = %q, want %q", tt.version, tt.scheme, got, tt.want)
			}
		})
	}
}

