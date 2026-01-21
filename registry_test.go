package purl

import (
	"testing"
)

func TestRegistryURL(t *testing.T) {
	tests := []struct {
		purl    string
		want    string
		wantErr bool
	}{
		// npm
		{"pkg:npm/lodash", "https://www.npmjs.com/package/lodash", false},
		{"pkg:npm/%40babel/core", "https://www.npmjs.com/package/@babel/core", false},

		// pypi
		{"pkg:pypi/requests", "https://pypi.org/project/requests/", false},

		// cargo
		{"pkg:cargo/serde", "https://crates.io/crates/serde", false},

		// gem
		{"pkg:gem/rails", "https://rubygems.org/gems/rails", false},

		// maven
		{"pkg:maven/org.apache.commons/commons-lang3", "https://mvnrepository.com/artifact/org.apache.commons/commons-lang3", false},

		// composer
		{"pkg:composer/symfony/console", "https://packagist.org/packages/symfony/console", false},

		// hex
		{"pkg:hex/phoenix", "https://hex.pm/packages/phoenix", false},

		// hackage
		{"pkg:hackage/aeson", "https://hackage.haskell.org/package/aeson", false},

		// pub
		{"pkg:pub/http", "https://pub.dev/packages/http", false},

		// nuget
		{"pkg:nuget/Newtonsoft.Json", "https://www.nuget.org/packages/Newtonsoft.Json", false},

		// conda
		{"pkg:conda/numpy", "https://anaconda.org/conda-forge/numpy", false},

		// homebrew
		{"pkg:homebrew/wget", "https://formulae.brew.sh/formula/wget", false},

		// deno
		{"pkg:deno/oak", "https://deno.land/x/oak", false},

		// No registry config
		{"pkg:deb/debian/curl", "", true},
		{"pkg:apk/alpine/curl", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.purl, func(t *testing.T) {
			p, err := Parse(tt.purl)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			got, err := p.RegistryURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("RegistryURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegistryURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRegistryURLWithVersion(t *testing.T) {
	tests := []struct {
		purl    string
		want    string
		wantErr bool
	}{
		// npm with version
		{"pkg:npm/lodash@4.17.21", "https://www.npmjs.com/package/lodash/v/4.17.21", false},
		{"pkg:npm/%40babel/core@7.24.0", "https://www.npmjs.com/package/@babel/core/v/7.24.0", false},

		// npm without version falls back
		{"pkg:npm/lodash", "https://www.npmjs.com/package/lodash", false},

		// pypi with version
		{"pkg:pypi/requests@2.28.1", "https://pypi.org/project/requests/2.28.1/", false},

		// gem with version
		{"pkg:gem/rails@7.0.4", "https://rubygems.org/gems/rails/versions/7.0.4", false},

		// nuget with version
		{"pkg:nuget/Newtonsoft.Json@13.0.1", "https://www.nuget.org/packages/Newtonsoft.Json/13.0.1", false},

		// maven with version
		{"pkg:maven/org.apache.commons/commons-lang3@3.12.0", "https://mvnrepository.com/artifact/org.apache.commons/commons-lang3/3.12.0", false},

		// hackage with version
		{"pkg:hackage/aeson@2.1.1.0", "https://hackage.haskell.org/package/aeson-2.1.1.0", false},

		// deno with version
		{"pkg:deno/oak@12.0.0", "https://deno.land/x/oak@12.0.0", false},

		// elm with version
		{"pkg:elm/elm/http@2.0.0", "https://package.elm-lang.org/packages/elm/http/2.0.0", false},

		// cargo without version support in URL
		{"pkg:cargo/serde@1.0.152", "https://crates.io/crates/serde", false},

		// hex without version support in URL
		{"pkg:hex/phoenix@1.7.0", "https://hex.pm/packages/phoenix", false},
	}

	for _, tt := range tests {
		t.Run(tt.purl, func(t *testing.T) {
			p, err := Parse(tt.purl)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			got, err := p.RegistryURLWithVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("RegistryURLWithVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegistryURLWithVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseRegistryURLWithType(t *testing.T) {
	tests := []struct {
		url       string
		purlType  string
		wantPURL  string
		wantErr   bool
	}{
		// npm
		{"https://www.npmjs.com/package/lodash", "npm", "pkg:npm/lodash", false},
		{"https://npmjs.com/package/lodash", "npm", "pkg:npm/lodash", false},
		{"https://www.npmjs.com/package/@babel/core", "npm", "pkg:npm/%40babel/core", false},
		{"https://www.npmjs.com/package/lodash/v/4.17.21", "npm", "pkg:npm/lodash@4.17.21", false},

		// pypi
		{"https://pypi.org/project/requests/", "pypi", "pkg:pypi/requests", false},
		{"https://pypi.org/project/requests/2.28.1/", "pypi", "pkg:pypi/requests@2.28.1", false},

		// cargo
		{"https://crates.io/crates/serde", "cargo", "pkg:cargo/serde", false},

		// gem
		{"https://rubygems.org/gems/rails", "gem", "pkg:gem/rails", false},
		{"https://rubygems.org/gems/rails/versions/7.0.4", "gem", "pkg:gem/rails@7.0.4", false},

		// maven
		{"https://mvnrepository.com/artifact/org.apache.commons/commons-lang3", "maven", "pkg:maven/org.apache.commons/commons-lang3", false},
		{"https://mvnrepository.com/artifact/org.apache.commons/commons-lang3/3.12.0", "maven", "pkg:maven/org.apache.commons/commons-lang3@3.12.0", false},

		// nuget
		{"https://www.nuget.org/packages/Newtonsoft.Json", "nuget", "pkg:nuget/Newtonsoft.Json", false},
		{"https://nuget.org/packages/Newtonsoft.Json/13.0.1", "nuget", "pkg:nuget/Newtonsoft.Json@13.0.1", false},

		// composer
		{"https://packagist.org/packages/symfony/console", "composer", "pkg:composer/symfony/console", false},

		// No match
		{"https://example.com/package", "npm", "", true},

		// No registry config
		{"https://example.com/package", "deb", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			p, err := ParseRegistryURLWithType(tt.url, tt.purlType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRegistryURLWithType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got := p.String(); got != tt.wantPURL {
				t.Errorf("ParseRegistryURLWithType() = %q, want %q", got, tt.wantPURL)
			}
		})
	}
}

func TestParseRegistryURL(t *testing.T) {
	tests := []struct {
		url      string
		wantType string
		wantErr  bool
	}{
		{"https://www.npmjs.com/package/lodash", "npm", false},
		{"https://pypi.org/project/requests/", "pypi", false},
		{"https://crates.io/crates/serde", "cargo", false},
		{"https://rubygems.org/gems/rails", "gem", false},
		{"https://example.com/unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			p, err := ParseRegistryURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRegistryURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if p.Type != tt.wantType {
				t.Errorf("ParseRegistryURL() type = %q, want %q", p.Type, tt.wantType)
			}
		})
	}
}
