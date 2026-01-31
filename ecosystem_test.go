package purl

import (
	"testing"
)

func TestNormalizeEcosystem(t *testing.T) {
	tests := []struct {
		ecosystem string
		want      string
	}{
		{"npm", "npm"},
		{"NPM", "npm"},
		{"go", "golang"},
		{"Go", "golang"},
		{"golang", "golang"},
		{"gem", "rubygems"},
		{"rubygems", "rubygems"},
		{"composer", "packagist"},
		{"packagist", "packagist"},
		{"cargo", "cargo"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.ecosystem, func(t *testing.T) {
			if got := NormalizeEcosystem(tt.ecosystem); got != tt.want {
				t.Errorf("NormalizeEcosystem(%q) = %q, want %q", tt.ecosystem, got, tt.want)
			}
		})
	}
}

func TestEcosystemToPURLType(t *testing.T) {
	tests := []struct {
		ecosystem string
		want      string
	}{
		{"npm", "npm"},
		{"alpine", "apk"},
		{"arch", "alpm"},
		{"rubygems", "gem"},
		{"packagist", "composer"},
		{"cargo", "cargo"},
		{"go", "golang"},           // alias normalized first
		{"gem", "gem"},             // alias to rubygems, then rubygems -> gem
		{"composer", "composer"},   // alias to packagist, then packagist -> composer
		{"github-actions", "githubactions"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.ecosystem, func(t *testing.T) {
			if got := EcosystemToPURLType(tt.ecosystem); got != tt.want {
				t.Errorf("EcosystemToPURLType(%q) = %q, want %q", tt.ecosystem, got, tt.want)
			}
		})
	}
}

func TestPURLTypeToEcosystem(t *testing.T) {
	tests := []struct {
		purlType string
		want     string
	}{
		{"apk", "alpine"},
		{"alpm", "arch"},
		{"gem", "rubygems"},
		{"composer", "packagist"},
		{"githubactions", "github-actions"},
		{"npm", "npm"},     // no reverse mapping, returns as-is
		{"cargo", "cargo"}, // no reverse mapping, returns as-is
	}

	for _, tt := range tests {
		t.Run(tt.purlType, func(t *testing.T) {
			if got := PURLTypeToEcosystem(tt.purlType); got != tt.want {
				t.Errorf("PURLTypeToEcosystem(%q) = %q, want %q", tt.purlType, got, tt.want)
			}
		})
	}
}

func TestEcosystemToOSV(t *testing.T) {
	tests := []struct {
		ecosystem string
		want      string
	}{
		{"npm", "npm"},
		{"rubygems", "RubyGems"},
		{"gem", "RubyGems"},
		{"pypi", "PyPI"},
		{"cargo", "crates.io"},
		{"golang", "Go"},
		{"go", "Go"},
		{"maven", "Maven"},
		{"nuget", "NuGet"},
		{"packagist", "Packagist"},
		{"composer", "Packagist"},
		{"hex", "Hex"},
		{"pub", "Pub"},
		{"cocoapods", "CocoaPods"},
		{"github-actions", "GitHub Actions"},
		{"unknown", "unknown"}, // falls through
	}

	for _, tt := range tests {
		t.Run(tt.ecosystem, func(t *testing.T) {
			if got := EcosystemToOSV(tt.ecosystem); got != tt.want {
				t.Errorf("EcosystemToOSV(%q) = %q, want %q", tt.ecosystem, got, tt.want)
			}
		})
	}
}

func TestMakePURL(t *testing.T) {
	tests := []struct {
		name      string
		ecosystem string
		pkg       string
		version   string
		wantStr   string
	}{
		// npm with scope
		{
			name:      "npm scoped",
			ecosystem: "npm",
			pkg:       "@babel/core",
			version:   "7.24.0",
			wantStr:   "pkg:npm/%40babel/core@7.24.0",
		},
		// npm without scope
		{
			name:      "npm unscoped",
			ecosystem: "npm",
			pkg:       "lodash",
			version:   "4.17.21",
			wantStr:   "pkg:npm/lodash@4.17.21",
		},
		// golang
		{
			name:      "golang",
			ecosystem: "golang",
			pkg:       "github.com/foo/bar",
			version:   "v1.0.0",
			wantStr:   "pkg:golang/github.com/foo/bar@v1.0.0",
		},
		// maven
		{
			name:      "maven",
			ecosystem: "maven",
			pkg:       "org.apache.commons:commons-lang3",
			version:   "3.12.0",
			wantStr:   "pkg:maven/org.apache.commons/commons-lang3@3.12.0",
		},
		// packagist/composer
		{
			name:      "packagist",
			ecosystem: "packagist",
			pkg:       "laravel/framework",
			version:   "9.0.0",
			wantStr:   "pkg:composer/laravel/framework@9.0.0",
		},
		// alpine (default namespace)
		{
			name:      "alpine",
			ecosystem: "alpine",
			pkg:       "curl",
			version:   "8.0.0",
			wantStr:   "pkg:apk/alpine/curl@8.0.0",
		},
		// arch (default namespace)
		{
			name:      "arch",
			ecosystem: "arch",
			pkg:       "base",
			version:   "1.0",
			wantStr:   "pkg:alpm/arch/base@1.0",
		},
		// cargo (no special handling)
		{
			name:      "cargo",
			ecosystem: "cargo",
			pkg:       "serde",
			version:   "1.0.0",
			wantStr:   "pkg:cargo/serde@1.0.0",
		},
		// pypi
		{
			name:      "pypi",
			ecosystem: "pypi",
			pkg:       "requests",
			version:   "2.28.0",
			wantStr:   "pkg:pypi/requests@2.28.0",
		},
		// using alias
		{
			name:      "go alias",
			ecosystem: "go",
			pkg:       "github.com/user/repo",
			version:   "v1.0.0",
			wantStr:   "pkg:golang/github.com/user/repo@v1.0.0",
		},
		// github-actions
		{
			name:      "github-actions",
			ecosystem: "github-actions",
			pkg:       "actions/checkout",
			version:   "v4",
			wantStr:   "pkg:githubactions/actions/checkout@v4",
		},
		// github-actions with path
		{
			name:      "github-actions with path",
			ecosystem: "github-actions",
			pkg:       "actions/cache/restore",
			version:   "v3",
			wantStr:   "pkg:githubactions/actions/cache@v3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := MakePURL(tt.ecosystem, tt.pkg, tt.version)
			if got := p.String(); got != tt.wantStr {
				t.Errorf("MakePURL(%q, %q, %q).String() = %q, want %q",
					tt.ecosystem, tt.pkg, tt.version, got, tt.wantStr)
			}
		})
	}
}

func TestMakePURLString(t *testing.T) {
	got := MakePURLString("npm", "lodash", "4.17.21")
	want := "pkg:npm/lodash@4.17.21"
	if got != want {
		t.Errorf("MakePURLString() = %q, want %q", got, want)
	}
}

func TestSupportedEcosystems(t *testing.T) {
	ecosystems := SupportedEcosystems()

	// Should have at least the known types plus aliases
	if len(ecosystems) < 10 {
		t.Errorf("SupportedEcosystems() returned only %d items, expected more", len(ecosystems))
	}

	// Check that some known ecosystems are present
	has := make(map[string]bool)
	for _, e := range ecosystems {
		has[e] = true
	}

	required := []string{"npm", "cargo", "pypi", "maven", "golang", "go", "gem", "alpine"}
	for _, r := range required {
		if !has[r] {
			t.Errorf("SupportedEcosystems() missing %q", r)
		}
	}
}

func TestIsValidEcosystem(t *testing.T) {
	tests := []struct {
		ecosystem string
		want      bool
	}{
		{"npm", true},
		{"cargo", true},
		{"pypi", true},
		{"golang", true},
		{"go", true},
		{"gem", true},
		{"rubygems", true},
		{"alpine", true},
		{"notarealecosystem", false},
	}

	for _, tt := range tests {
		t.Run(tt.ecosystem, func(t *testing.T) {
			if got := IsValidEcosystem(tt.ecosystem); got != tt.want {
				t.Errorf("IsValidEcosystem(%q) = %v, want %v", tt.ecosystem, got, tt.want)
			}
		})
	}
}
