package purl

import "testing"

func TestComparisonForm(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// PEP 503: case and runs of -, _ and . are equivalent.
		{"pkg:pypi/zope.interface@6.0", "pkg:pypi/zope-interface@6.0"},
		{"pkg:pypi/zope_interface@6.0", "pkg:pypi/zope-interface@6.0"},
		{"pkg:pypi/Zope-Interface@6.0", "pkg:pypi/zope-interface@6.0"},
		{"pkg:pypi/zope-interface@6.0", "pkg:pypi/zope-interface@6.0"},
		{"pkg:pypi/Foo._-Bar", "pkg:pypi/foo-bar"},
		{"pkg:pypi/Django@4.2RC1", "pkg:pypi/django@4.2RC1"},

		// NuGet ids and versions are case-insensitive.
		{"pkg:nuget/Newtonsoft.Json@13.0.1", "pkg:nuget/newtonsoft.json@13.0.1"},
		{"pkg:nuget/newtonsoft.json@13.0.1", "pkg:nuget/newtonsoft.json@13.0.1"},
		{"pkg:nuget/NEWTONSOFT.JSON@13.0.1", "pkg:nuget/newtonsoft.json@13.0.1"},
		{"pkg:nuget/Foo.Bar@1.0.0-Beta.1", "pkg:nuget/foo.bar@1.0.0-beta.1"},

		// Composer drops a tag's leading v.
		{"pkg:composer/symfony/console@v6.1.7", "pkg:composer/symfony/console@6.1.7"},
		{"pkg:composer/symfony/console@V6.1.7", "pkg:composer/symfony/console@6.1.7"},
		{"pkg:composer/symfony/console@6.1.7", "pkg:composer/symfony/console@6.1.7"},
		{"pkg:composer/foo/bar@v", "pkg:composer/foo/bar@v"},
		{"pkg:composer/foo/bar@vendor-branch", "pkg:composer/foo/bar@vendor-branch"},
		{"pkg:composer/foo/bar@dev-main", "pkg:composer/foo/bar@dev-main"},

		// conda: the channel ends up in the qualifier however it was given.
		{"pkg:conda/numpy@1.26.0", "pkg:conda/numpy@1.26.0"},
		{"pkg:conda/numpy@1.26.0?channel=conda-forge", "pkg:conda/numpy@1.26.0?channel=conda-forge"},
		{"pkg:conda/conda-forge/numpy@1.26.0", "pkg:conda/numpy@1.26.0?channel=conda-forge"},
		{"pkg:conda/conda-forge%2Fnumpy@1.26.0", "pkg:conda/numpy@1.26.0?channel=conda-forge"},
		{"pkg:conda/conda-forge/numpy@1.26.0?channel=main", "pkg:conda/numpy@1.26.0?channel=conda-forge"},
		{"pkg:conda/numpy@1.26.0?build=py312&channel=main", "pkg:conda/numpy@1.26.0?build=py312&channel=main"},

		// oci: a path-style repository is one name.
		{"pkg:oci/library/alpine@3.19", "pkg:oci/library%2Falpine@3.19"},
		{"pkg:oci/library%2Falpine@3.19", "pkg:oci/library%2Falpine@3.19"},
		{"pkg:oci/alpine@3.19", "pkg:oci/alpine@3.19"},

		// swift: registry identities are case-insensitive, source coordinates are not.
		{"pkg:swift/Apple/Swift-NIO@2.0.0", "pkg:swift/apple/swift-nio@2.0.0"},
		{"pkg:swift/github.com/Apple/Swift-NIO@2.0.0", "pkg:swift/github.com/Apple/Swift-NIO@2.0.0"},

		// Types without extra rules are unchanged.
		{"pkg:npm/%40babel/core@7.24.0", "pkg:npm/%40babel/core@7.24.0"},
		{"pkg:gem/Rails@7.0.0", "pkg:gem/Rails@7.0.0"},
		{"pkg:maven/org.Apache/Commons@1.0", "pkg:maven/org.Apache/Commons@1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.input, err)
			}
			if got := p.ComparisonForm().String(); got != tt.want {
				t.Errorf("ComparisonForm(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestComparisonFormEquivalentSpellings(t *testing.T) {
	groups := [][]string{
		{"pkg:pypi/zope.interface@6.0", "pkg:pypi/zope_interface@6.0", "pkg:pypi/Zope-Interface@6.0", "pkg:pypi/ZOPE__INTERFACE@6.0"},
		{"pkg:nuget/Newtonsoft.Json@13.0.1-Beta", "pkg:nuget/newtonsoft.json@13.0.1-beta"},
		{"pkg:composer/symfony/console@v6.1.7", "pkg:composer/symfony/console@6.1.7"},
		{"pkg:conda/conda-forge/numpy@1.26.0", "pkg:conda/numpy@1.26.0?channel=conda-forge", "pkg:conda/conda-forge%2Fnumpy@1.26.0"},
		{"pkg:oci/library/alpine@3.19", "pkg:oci/library%2Falpine@3.19"},
		{"pkg:swift/Apple/Swift-NIO@2.0.0", "pkg:swift/apple/swift-nio@2.0.0"},
	}
	for _, group := range groups {
		var want string
		for i, s := range group {
			p, err := Parse(s)
			if err != nil {
				t.Fatalf("Parse(%q): %v", s, err)
			}
			got := p.ComparisonForm().String()
			if i == 0 {
				want = got
			} else if got != want {
				t.Errorf("ComparisonForm(%q) = %q, want %q (as %q)", s, got, want, group[0])
			}
		}
	}
}

// Parse and String feed storage keys, so ComparisonForm must not change them.
func TestComparisonFormLeavesParseAlone(t *testing.T) {
	inputs := map[string]string{
		"pkg:pypi/zope.interface@6.0":         "pkg:pypi/zope.interface@6.0",
		"pkg:pypi/Zope_Interface@6.0":         "pkg:pypi/zope-interface@6.0",
		"pkg:nuget/Newtonsoft.Json@13.0.1":    "pkg:nuget/Newtonsoft.Json@13.0.1",
		"pkg:composer/symfony/console@v6.1.7": "pkg:composer/symfony/console@v6.1.7",
		"pkg:conda/conda-forge/numpy@1.26.0":  "pkg:conda/conda-forge/numpy@1.26.0",
		"pkg:oci/library/alpine@3.19":         "pkg:oci/library/alpine@3.19",
		"pkg:swift/Apple/Swift-NIO@2.0.0":     "pkg:swift/Apple/Swift-NIO@2.0.0",
	}
	for input, want := range inputs {
		p, err := Parse(input)
		if err != nil {
			t.Fatalf("Parse(%q): %v", input, err)
		}
		_ = p.ComparisonForm()
		if got := p.String(); got != want {
			t.Errorf("Parse(%q).String() = %q after ComparisonForm, want %q", input, got, want)
		}
	}
}

func TestComparisonFormDoesNotShareQualifiers(t *testing.T) {
	p, err := Parse("pkg:conda/conda-forge/numpy@1.26.0?channel=main")
	if err != nil {
		t.Fatal(err)
	}
	_ = p.ComparisonForm()
	if got := p.Qualifier("channel"); got != "main" {
		t.Errorf("original channel = %q, want main", got)
	}
}
