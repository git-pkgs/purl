package purl_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/git-pkgs/purl"
)

func TestSpecTypeMetadata(t *testing.T) {
	var types []string
	readSpecJSON(t, "purl-types-index.json", &types)
	if len(types) == 0 {
		t.Fatal("spec type index is empty")
	}
	registered := make(map[string]bool, len(types))
	for _, name := range types {
		if registered[name] {
			t.Errorf("duplicate spec type %q", name)
		}
		registered[name] = true
		t.Run(name, func(t *testing.T) {
			checkSpecTypeMetadata(t, name)
		})
	}

	localTypes := map[string]string{
		"bioconductor": "Bioconductor registry support",
		"bsd":          "FreeBSD packages, proposed in purl-spec#496",
		"clojars":      "Clojars registry support",
		"deno":         "Deno registry support",
		"elm":          "Elm registry support",
		"homebrew":     "Legacy type retained alongside brew for compatibility",
		"terraform":    "Terraform providers and modules, proposed in purl-spec#889",
	}
	for _, name := range purl.KnownTypes() {
		if !registered[name] && localTypes[name] == "" {
			t.Errorf("unregistered type %q needs an explicit local type allowance", name)
		}
	}
	for name := range localTypes {
		if registered[name] || !purl.IsKnownType(name) {
			t.Errorf("remove stale local type allowance for %q", name)
		}
	}
}

func checkSpecTypeMetadata(t *testing.T, name string) {
	t.Helper()
	var definition struct {
		Type      string `json:"type"`
		Namespace struct {
			Requirement string `json:"requirement"`
		} `json:"namespace_definition"`
	}
	readSpecJSON(t, filepath.Join("types", name+"-definition.json"), &definition)
	if definition.Type != name {
		t.Fatalf("definition type = %q, want %q", definition.Type, name)
	}
	requirement := definition.Namespace.Requirement
	switch requirement {
	case "required", "optional", "prohibited":
	default:
		t.Fatalf("unknown namespace requirement %q", requirement)
	}
	cfg := purl.TypeInfo(name)
	if cfg == nil {
		t.Fatalf("spec type %q is missing from types.json", name)
	}
	if !purl.IsKnownType(name) || !purl.IsValidEcosystem(name) {
		t.Errorf("spec type %q is not recognized", name)
	}
	if cfg.NamespaceRequirement != requirement {
		t.Errorf("namespace requirement = %q, spec requires %q", cfg.NamespaceRequirement, requirement)
	}
	if cfg.NamespaceRequired() != (requirement == "required") || cfg.NamespaceProhibited() != (requirement == "prohibited") {
		t.Errorf("namespace helpers disagree with spec requirement %q", requirement)
	}
}

func readSpecJSON(t *testing.T, name string, value any) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "purl-spec", name))
	if err != nil {
		t.Fatalf("read spec: %v; initialize it with git submodule update --init", err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		t.Fatalf("decode spec %s: %v", name, err)
	}
}

func TestHomebrewCompatibility(t *testing.T) {
	const registry = "https://formulae.brew.sh"
	const packageURL = registry + "/formula/wget"
	const want = "pkg:homebrew/wget@1.21.3"
	for _, got := range []string{
		purl.MakePURL("homebrew", "wget", "1.21.3").String(),
		purl.MakePURLString("homebrew", "wget", "1.21.3"),
		purl.BuildPURLString("homebrew", "wget", "1.21.3", registry),
	} {
		if got != want {
			t.Errorf("homebrew builder = %q, want %q", got, want)
		}
	}
	p, err := purl.Parse(want)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := p.RegistryURL(); err != nil || got != packageURL {
		t.Errorf("RegistryURL() = %q, %v", got, err)
	}
	p, err = purl.ParseRegistryURL(packageURL)
	if err != nil {
		t.Fatal(err)
	}
	if got := p.String(); got != "pkg:homebrew/wget" {
		t.Errorf("ParseRegistryURL() = %q, want legacy homebrew type", got)
	}
	if !purl.IsKnownType("brew") || !purl.IsKnownType("homebrew") {
		t.Error("both brew and homebrew must remain recognized")
	}
	if got := purl.MakePURLString("brew", "wget", "1.21.3"); got != "pkg:brew/wget@1.21.3" {
		t.Errorf("explicit brew builder = %q", got)
	}
}
