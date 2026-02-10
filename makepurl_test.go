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

func TestBuildPURLString(t *testing.T) {
	tests := []struct {
		name        string
		ecosystem   string
		pkgName     string
		version     string
		registryURL string
		want        string
	}{
		{"simple npm", "npm", "lodash", "4.17.21", "", "pkg:npm/lodash@4.17.21"},
		{"scoped npm", "npm", "@babel/core", "7.20.0", "", "pkg:npm/%40babel/core@7.20.0"},  // @ encoded in namespace
		{"gem", "rubygems", "rails", "7.0.0", "", "pkg:gem/rails@7.0.0"},
		{"pypi", "pypi", "requests", "2.28.0", "", "pkg:pypi/requests@2.28.0"},
		{"maven", "maven", "org.apache:commons", "1.0", "", "pkg:maven/org.apache/commons@1.0"},
		{"golang", "golang", "github.com/foo/bar", "v1.0.0", "", "pkg:golang/github.com/foo/bar@v1.0.0"},
		{"no version", "npm", "lodash", "", "", "pkg:npm/lodash"},
		{"with registry", "npm", "lodash", "1.0.0", "https://npm.example.com", "pkg:npm/lodash@1.0.0?repository_url=https:%2F%2Fnpm.example.com"},
		{"default registry ignored", "npm", "lodash", "1.0.0", "https://registry.npmjs.org", "pkg:npm/lodash@1.0.0"},
		{"composer", "packagist", "vendor/pkg", "1.0", "", "pkg:composer/vendor/pkg@1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPURLString(tt.ecosystem, tt.pkgName, tt.version, tt.registryURL)
			if got != tt.want {
				t.Errorf("BuildPURLString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildPURLStringMatchesMakePURL(t *testing.T) {
	// Ensure BuildPURLString produces the same output as the struct-based path
	cases := []struct {
		ecosystem   string
		name        string
		version     string
		registryURL string
	}{
		{"npm", "lodash", "4.17.21", ""},
		{"npm", "@babel/core", "7.20.0", ""},
		{"rubygems", "rails", "7.0.0", ""},
		{"maven", "org.apache:commons", "1.0", ""},
		{"golang", "github.com/foo/bar", "v1.0.0", ""},
		{"packagist", "vendor/pkg", "1.0", ""},
		{"npm", "lodash", "^1.0.0", ""},
		{"npm", "pkg", "1.0.0", "https://custom.registry.com"},
	}

	for _, tt := range cases {
		t.Run(tt.ecosystem+"/"+tt.name, func(t *testing.T) {
			fast := BuildPURLString(tt.ecosystem, tt.name, tt.version, tt.registryURL)

			purlType := EcosystemToPURLType(tt.ecosystem)
			cleanVersion := CleanVersion(tt.version, purlType)
			p := MakePURL(tt.ecosystem, tt.name, cleanVersion)
			if tt.registryURL != "" && IsNonDefaultRegistry(purlType, tt.registryURL) {
				p = p.WithQualifier("repository_url", tt.registryURL)
			}
			slow := p.String()

			if fast != slow {
				t.Errorf("mismatch:\n  fast: %q\n  slow: %q", fast, slow)
			}
		})
	}
}

