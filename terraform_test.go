package purl_test

import (
	"testing"

	"github.com/git-pkgs/purl"
)

// Terraform folds the case of a provider's namespace and type, and keeps a
// module's as given.
func TestTerraformProviderIsLowercased(t *testing.T) {
	const want = "pkg:terraform/hashicorp/random@3.6.0"
	parsed, err := purl.Parse("pkg:terraform/HashiCorp/Random@3.6.0")
	if err != nil {
		t.Fatal(err)
	}
	for _, got := range []string{
		parsed.String(),
		purl.New("terraform", "HashiCorp", "Random", "3.6.0", nil).String(),
		purl.MakePURL("terraform", "HashiCorp/Random", "3.6.0").String(),
		purl.MakePURLString("terraform", "HashiCorp/Random", "3.6.0"),
		purl.BuildPURLString("terraform", "HashiCorp/Random", "3.6.0", ""),
	} {
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestTerraformModuleKeepsCase(t *testing.T) {
	const want = "pkg:terraform/HashiCorp/Consul/AWS@0.11.0"
	parsed, err := purl.Parse(want)
	if err != nil {
		t.Fatal(err)
	}
	for _, got := range []string{
		parsed.String(),
		purl.New("terraform", "HashiCorp/Consul", "AWS", "0.11.0", nil).String(),
		purl.MakePURL("terraform", "HashiCorp/Consul/AWS", "0.11.0").String(),
		purl.MakePURLString("terraform", "HashiCorp/Consul/AWS", "0.11.0"),
		purl.BuildPURLString("terraform", "HashiCorp/Consul/AWS", "0.11.0", ""),
	} {
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestTerraformRejectsHostnamesAndExtraSegments(t *testing.T) {
	for _, name := range []string{
		"registry.terraform.io/hashicorp/random",
		"registry.terraform.io/hashicorp/consul/aws",
		"example.com/random",
		"a/b/c/d",
	} {
		if got := purl.MakePURLString("terraform", name, "1.0.0"); got != "" {
			t.Errorf("MakePURLString(%q) = %q, want empty", name, got)
		}
		if got := purl.MakePURL("terraform", name, "1.0.0"); got != nil {
			t.Errorf("MakePURL(%q) = %v, want nil", name, got)
		}
	}
}
