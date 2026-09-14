package purl

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildPURLRegistryHostSemantics(t *testing.T) {
	for _, tc := range []struct {
		registry  string
		isDefault bool
	}{
		{"https://registry.npmjs.org/package/-/package.tgz", true},
		{"http://registry.npmjs.org:8080/package", true},
		{"//registry.npmjs.org/package", true},
		{"https://cdn.registry.npmjs.org/package", true},
		{"https://user:pass@registry.npmjs.org/package", true},
		{"https://evil.example@registry.npmjs.org/package", true},
		{"https://registry.npmjs.org@evil.example/package", false},
		{"https://registry.npmjs.org.evil.example/package", false},
		{"https://notregistry.npmjs.org/package", false},
		{"https://REGISTRY.NPMJS.ORG/package", false},
		{"https://registry.npmjs.org./package", false},
		{"https://registry.npmjs.org:bad/package", false},
		{"https://registry.npmjs.org/%zz", false},
		{"https://[::1]/package", false},
		{"registry.npmjs.org/package", false},
	} {
		t.Run(tc.registry, func(t *testing.T) {
			got := BuildPURLString("npm", "package", "1.2.3", tc.registry)
			if tc.isDefault {
				if got != "pkg:npm/package@1.2.3" {
					t.Fatalf("default registry retained: %q", got)
				}
			} else if !strings.HasPrefix(got, "pkg:npm/package@1.2.3?repository_url=") {
				t.Fatalf("registry lost: %q", got)
			}
		})
	}
}

func TestDefaultRegistryPointerChanges(t *testing.T) {
	cfg := TypeInfo("npm")
	original := *cfg.DefaultRegistry
	t.Cleanup(func() { *cfg.DefaultRegistry = original })
	for _, registry := range []string{"https://custom.example/registry", "https://crates.io", "", "https://[invalid", original} {
		*cfg.DefaultRegistry = registry
		for _, input := range []string{original, "https://custom.example/package.tgz", "https://crates.io", ""} {
			want := uncachedDefaultRegistry("npm", input)
			if got := IsDefaultRegistry("npm", input); got != want {
				t.Fatalf("default %q, input %q: got %t, want %t", registry, input, got, want)
			}
			got := BuildPURLString("npm", "package", "1.2.3", input)
			if strings.Contains(got, "?repository_url=") == want {
				t.Fatalf("default %q, input %q: unexpected PURL %q", registry, input, got)
			}
		}
	}
}

func uncachedDefaultRegistry(purlType, registry string) bool {
	if registry == "" {
		return true
	}
	cfg := TypeInfo(purlType)
	if cfg == nil || cfg.DefaultRegistry == nil {
		return false
	}
	standard, err := url.Parse(*cfg.DefaultRegistry)
	if err != nil || standard.Hostname() == "" {
		return false
	}
	given, err := url.Parse(registry)
	if err != nil || given.Hostname() == "" {
		return false
	}
	return given.Hostname() == standard.Hostname() || strings.HasSuffix(given.Hostname(), "."+standard.Hostname())
}

func FuzzDefaultRegistryCompatibility(f *testing.F) {
	for _, kind := range append(KnownTypes(), "unknown") {
		f.Add(kind, DefaultRegistry(kind))
		f.Add(kind, "https://user@registry.npmjs.org:443/package.tgz")
		f.Add(kind, "https://registry.npmjs.org@evil.example/%zz")
	}
	f.Fuzz(func(t *testing.T, kind, registry string) {
		want := uncachedDefaultRegistry(kind, registry)
		if got := IsDefaultRegistry(kind, registry); got != want {
			t.Fatalf("IsDefaultRegistry(%q, %q) = %t, want %t", kind, registry, got, want)
		}
		if got := IsNonDefaultRegistry(kind, registry); got == want {
			t.Fatalf("IsNonDefaultRegistry(%q, %q) = %t, want %t", kind, registry, got, !want)
		}
	})
}
