package purl

import (
	"testing"
)

func TestIsDefaultRegistry(t *testing.T) {
	tests := []struct {
		purlType    string
		registryURL string
		want        bool
	}{
		// Empty URL is always default
		{"npm", "", true},
		{"pypi", "", true},

		// npm default registry (from types.json: https://registry.npmjs.org)
		{"npm", "https://registry.npmjs.org", true},
		{"npm", "https://npm.example.com", false},

		// pypi default registry (from types.json: https://pypi.org)
		{"pypi", "https://pypi.org", true},
		{"pypi", "https://pypi.example.com", false},

		// cargo default registry (from types.json: https://crates.io)
		{"cargo", "https://crates.io", true},
		{"cargo", "https://cargo.example.com", false},

		// gem default registry (from types.json: https://rubygems.org)
		{"gem", "https://rubygems.org", true},
		{"gem", "https://gems.example.com", false},

		// maven default registry (from types.json: https://repo.maven.apache.org/maven2)
		{"maven", "https://repo.maven.apache.org", true},
		{"maven", "https://maven.example.com", false},

		// golang default registry (from types.json: https://pkg.go.dev)
		{"golang", "https://pkg.go.dev", true},
		{"golang", "https://go.example.com", false},

		// docker default registry (from types.json: https://hub.docker.com)
		{"docker", "https://hub.docker.com", true},
		{"docker", "https://gcr.io", false},

		// subdomain matching
		{"npm", "https://cdn.registry.npmjs.org", true},

		// Types with no default registry
		{"apk", "https://example.com", false},
		{"deb", "https://example.com", false},

		// Unknown type
		{"unknown-type", "https://example.com", false},
	}

	for _, tt := range tests {
		name := tt.purlType + ":" + tt.registryURL
		if tt.registryURL == "" {
			name = tt.purlType + ":empty"
		}
		t.Run(name, func(t *testing.T) {
			if got := IsDefaultRegistry(tt.purlType, tt.registryURL); got != tt.want {
				t.Errorf("IsDefaultRegistry(%q, %q) = %v, want %v", tt.purlType, tt.registryURL, got, tt.want)
			}
		})
	}
}

func TestIsNonDefaultRegistry(t *testing.T) {
	tests := []struct {
		purlType    string
		registryURL string
		want        bool
	}{
		{"npm", "", false},
		{"npm", "https://registry.npmjs.org", false},
		{"npm", "https://npm.example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.purlType+":"+tt.registryURL, func(t *testing.T) {
			if got := IsNonDefaultRegistry(tt.purlType, tt.registryURL); got != tt.want {
				t.Errorf("IsNonDefaultRegistry(%q, %q) = %v, want %v", tt.purlType, tt.registryURL, got, tt.want)
			}
		})
	}
}

func TestIsPrivateRegistry(t *testing.T) {
	tests := []struct {
		purl string
		want bool
	}{
		{"pkg:npm/lodash", false},
		{"pkg:npm/lodash?repository_url=https://registry.npmjs.org", false},
		{"pkg:npm/lodash?repository_url=https://npm.example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.purl, func(t *testing.T) {
			p, err := Parse(tt.purl)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if got := p.IsPrivateRegistry(); got != tt.want {
				t.Errorf("IsPrivateRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}
