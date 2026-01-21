package purl

import (
	"testing"
)

func TestTypeInfo(t *testing.T) {
	tests := []struct {
		purlType        string
		wantDescription string
		wantRegistry    string
		wantNil         bool
	}{
		{"npm", "PURL type for npm packages.", "https://registry.npmjs.org", false},
		{"pypi", "Python packages", "https://pypi.org", false},
		{"cargo", "Cargo packages for Rust", "https://crates.io", false},
		{"unknown-type", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.purlType, func(t *testing.T) {
			cfg := TypeInfo(tt.purlType)
			if tt.wantNil {
				if cfg != nil {
					t.Errorf("TypeInfo(%q) = %v, want nil", tt.purlType, cfg)
				}
				return
			}
			if cfg == nil {
				t.Fatalf("TypeInfo(%q) = nil, want non-nil", tt.purlType)
			}
			if cfg.Description != tt.wantDescription {
				t.Errorf("Description = %q, want %q", cfg.Description, tt.wantDescription)
			}
			if cfg.DefaultRegistry == nil {
				if tt.wantRegistry != "" {
					t.Errorf("DefaultRegistry = nil, want %q", tt.wantRegistry)
				}
			} else if *cfg.DefaultRegistry != tt.wantRegistry {
				t.Errorf("DefaultRegistry = %q, want %q", *cfg.DefaultRegistry, tt.wantRegistry)
			}
		})
	}
}

func TestKnownTypes(t *testing.T) {
	types := KnownTypes()
	if len(types) == 0 {
		t.Fatal("KnownTypes() returned empty slice")
	}

	// Check that some expected types are present
	expected := []string{"npm", "pypi", "cargo", "gem", "maven", "golang"}
	for _, e := range expected {
		found := false
		for _, typ := range types {
			if typ == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("KnownTypes() missing %q", e)
		}
	}

	// Verify sorted
	for i := 1; i < len(types); i++ {
		if types[i-1] > types[i] {
			t.Errorf("KnownTypes() not sorted: %q > %q", types[i-1], types[i])
			break
		}
	}
}

func TestIsKnownType(t *testing.T) {
	tests := []struct {
		purlType string
		want     bool
	}{
		{"npm", true},
		{"pypi", true},
		{"cargo", true},
		{"unknown-type", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.purlType, func(t *testing.T) {
			if got := IsKnownType(tt.purlType); got != tt.want {
				t.Errorf("IsKnownType(%q) = %v, want %v", tt.purlType, got, tt.want)
			}
		})
	}
}

func TestDefaultRegistry(t *testing.T) {
	tests := []struct {
		purlType string
		want     string
	}{
		{"npm", "https://registry.npmjs.org"},
		{"pypi", "https://pypi.org"},
		{"cargo", "https://crates.io"},
		{"apk", ""},  // no default registry
		{"deb", ""},  // no default registry
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.purlType, func(t *testing.T) {
			if got := DefaultRegistry(tt.purlType); got != tt.want {
				t.Errorf("DefaultRegistry(%q) = %q, want %q", tt.purlType, got, tt.want)
			}
		})
	}
}

func TestNamespaceRequirement(t *testing.T) {
	tests := []struct {
		purlType         string
		wantRequired     bool
		wantProhibited   bool
	}{
		{"maven", true, false},
		{"composer", true, false},
		{"gem", false, true},
		{"cran", false, true},
		{"npm", false, false},  // optional
		{"cargo", false, false}, // no requirement
	}

	for _, tt := range tests {
		t.Run(tt.purlType, func(t *testing.T) {
			cfg := TypeInfo(tt.purlType)
			if cfg == nil {
				t.Fatalf("TypeInfo(%q) = nil", tt.purlType)
			}
			if got := cfg.NamespaceRequired(); got != tt.wantRequired {
				t.Errorf("NamespaceRequired() = %v, want %v", got, tt.wantRequired)
			}
			if got := cfg.NamespaceProhibited(); got != tt.wantProhibited {
				t.Errorf("NamespaceProhibited() = %v, want %v", got, tt.wantProhibited)
			}
		})
	}
}

func TestTypesVersion(t *testing.T) {
	v := TypesVersion()
	if v == "" {
		t.Error("TypesVersion() returned empty string")
	}
}
