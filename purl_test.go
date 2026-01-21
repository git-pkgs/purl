package purl

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
		wantNS   string
		wantName string
		wantVer  string
		wantErr  bool
	}{
		// Basic packages without version
		{"pkg:cargo/serde", "cargo", "", "serde", "", false},
		{"pkg:npm/lodash", "npm", "", "lodash", "", false},
		{"pkg:pypi/requests", "pypi", "", "requests", "", false},
		{"pkg:gem/rails", "gem", "", "rails", "", false},

		// Packages with version
		{"pkg:cargo/serde@1.0.0", "cargo", "", "serde", "1.0.0", false},
		{"pkg:npm/lodash@4.17.21", "npm", "", "lodash", "4.17.21", false},
		{"pkg:gem/rails@7.0.0", "gem", "", "rails", "7.0.0", false},

		// npm scoped packages
		{"pkg:npm/%40babel/core", "npm", "@babel", "core", "", false},
		{"pkg:npm/%40babel/core@7.24.0", "npm", "@babel", "core", "7.24.0", false},

		// Maven with groupId
		{"pkg:maven/org.apache.commons/commons-lang3", "maven", "org.apache.commons", "commons-lang3", "", false},
		{"pkg:maven/org.apache.commons/commons-lang3@3.12.0", "maven", "org.apache.commons", "commons-lang3", "3.12.0", false},

		// Go modules
		{"pkg:golang/github.com/gorilla/mux", "golang", "github.com/gorilla", "mux", "", false},
		{"pkg:golang/github.com/gorilla/mux@v1.8.0", "golang", "github.com/gorilla", "mux", "v1.8.0", false},

		// Terraform modules
		{"pkg:terraform/hashicorp/consul/aws", "terraform", "hashicorp/consul", "aws", "", false},
		{"pkg:terraform/hashicorp/consul/aws@0.11.0", "terraform", "hashicorp/consul", "aws", "0.11.0", false},

		// Hex
		{"pkg:hex/phoenix@1.7.0", "hex", "", "phoenix", "1.7.0", false},

		// Composer
		{"pkg:composer/symfony/console@6.1.7", "composer", "symfony", "console", "6.1.7", false},

		// Errors
		{"cargo/serde", "", "", "", "", true}, // missing pkg: prefix
		{"invalid", "", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if p.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", p.Type, tt.wantType)
			}
			if p.Namespace != tt.wantNS {
				t.Errorf("Namespace = %q, want %q", p.Namespace, tt.wantNS)
			}
			if p.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", p.Name, tt.wantName)
			}
			if p.Version != tt.wantVer {
				t.Errorf("Version = %q, want %q", p.Version, tt.wantVer)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		purlType   string
		namespace  string
		pkgName    string
		version    string
		qualifiers map[string]string
		want       string
	}{
		{
			name:     "simple package",
			purlType: "cargo",
			pkgName:  "serde",
			want:     "pkg:cargo/serde",
		},
		{
			name:     "package with version",
			purlType: "npm",
			pkgName:  "lodash",
			version:  "4.17.21",
			want:     "pkg:npm/lodash@4.17.21",
		},
		{
			name:      "npm scoped",
			purlType:  "npm",
			namespace: "@babel",
			pkgName:   "core",
			version:   "7.24.0",
			want:      "pkg:npm/%40babel/core@7.24.0",
		},
		{
			name:      "maven",
			purlType:  "maven",
			namespace: "org.apache.commons",
			pkgName:   "commons-lang3",
			version:   "3.12.0",
			want:      "pkg:maven/org.apache.commons/commons-lang3@3.12.0",
		},
		{
			name:       "with qualifier",
			purlType:   "npm",
			pkgName:    "lodash",
			qualifiers: map[string]string{"repository_url": "https://npm.example.com"},
			want:       "pkg:npm/lodash?repository_url=https%3A%2F%2Fnpm.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.purlType, tt.namespace, tt.pkgName, tt.version, tt.qualifiers)
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRepositoryURL(t *testing.T) {
	tests := []struct {
		purl string
		want string
	}{
		{"pkg:npm/lodash", ""},
		{"pkg:npm/lodash?repository_url=https://npm.example.com", "https://npm.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.purl, func(t *testing.T) {
			p, err := Parse(tt.purl)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if got := p.RepositoryURL(); got != tt.want {
				t.Errorf("RepositoryURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWithVersion(t *testing.T) {
	p, _ := Parse("pkg:npm/lodash")
	pv := p.WithVersion("4.17.21")

	if p.Version != "" {
		t.Error("Original PURL should not have version")
	}
	if pv.Version != "4.17.21" {
		t.Errorf("WithVersion() version = %q, want %q", pv.Version, "4.17.21")
	}
	if pv.String() != "pkg:npm/lodash@4.17.21" {
		t.Errorf("WithVersion() string = %q, want %q", pv.String(), "pkg:npm/lodash@4.17.21")
	}
}

func TestWithoutVersion(t *testing.T) {
	p, _ := Parse("pkg:npm/lodash@4.17.21")
	pv := p.WithoutVersion()

	if p.Version != "4.17.21" {
		t.Error("Original PURL should have version")
	}
	if pv.Version != "" {
		t.Errorf("WithoutVersion() version = %q, want empty", pv.Version)
	}
}

func TestQualifier(t *testing.T) {
	p, _ := Parse("pkg:npm/lodash?repository_url=https://npm.example.com&checksum=sha256:abc123")

	tests := []struct {
		key  string
		want string
	}{
		{"repository_url", "https://npm.example.com"},
		{"checksum", "sha256:abc123"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := p.Qualifier(tt.key); got != tt.want {
				t.Errorf("Qualifier(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestWithQualifier(t *testing.T) {
	p, _ := Parse("pkg:npm/lodash@4.17.21")

	// Add a new qualifier
	p2 := p.WithQualifier("repository_url", "https://npm.example.com")
	if got := p2.Qualifier("repository_url"); got != "https://npm.example.com" {
		t.Errorf("WithQualifier() new qualifier = %q, want %q", got, "https://npm.example.com")
	}

	// Original should not be modified
	if got := p.Qualifier("repository_url"); got != "" {
		t.Errorf("Original should not have qualifier, got %q", got)
	}

	// Replace an existing qualifier
	p3 := p2.WithQualifier("repository_url", "https://other.example.com")
	if got := p3.Qualifier("repository_url"); got != "https://other.example.com" {
		t.Errorf("WithQualifier() replaced qualifier = %q, want %q", got, "https://other.example.com")
	}

	// Add a second qualifier
	p4 := p2.WithQualifier("checksum", "sha256:abc")
	if got := p4.Qualifier("checksum"); got != "sha256:abc" {
		t.Errorf("WithQualifier() second qualifier = %q, want %q", got, "sha256:abc")
	}
	if got := p4.Qualifier("repository_url"); got != "https://npm.example.com" {
		t.Errorf("WithQualifier() should preserve existing qualifier = %q, want %q", got, "https://npm.example.com")
	}
}
