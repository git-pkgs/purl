package purl_test

import (
	"testing"

	"github.com/git-pkgs/purl"
)

func TestVSCodeExtensionCase(t *testing.T) {
	const want = "pkg:vscode-extension/ms-python/python@2023.25.10292213"
	parsed, err := purl.Parse("pkg:vscode-extension/MS-Python/Python@2023.25.10292213")
	if err != nil {
		t.Fatal(err)
	}
	for _, got := range []string{
		parsed.String(),
		purl.New("vscode-extension", "MS-Python", "Python", "2023.25.10292213", nil).String(),
		purl.MakePURL("vscode-extension", "MS-Python.Python", "2023.25.10292213").String(),
		purl.MakePURLString("vscode-extension", "MS-Python.Python", "2023.25.10292213"),
		purl.BuildPURLString("vscode-extension", "MS-Python.Python", "2023.25.10292213", ""),
	} {
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestVSCodeExtensionOpenVSX(t *testing.T) {
	const want = "pkg:vscode-extension/muhammad-sammy/csharp@2.15.30-RC1?platform=linux-x64&repository_url=https:%2F%2Fopen-vsx.org"
	parsed, err := purl.Parse("pkg:vscode-extension/Muhammad-Sammy/CSharp@2.15.30-RC1?platform=linux-x64&repository_url=https:%2F%2Fopen-vsx.org")
	if err != nil {
		t.Fatal(err)
	}
	if got := parsed.String(); got != want {
		t.Errorf("Parse() = %q, want %q", got, want)
	}
	if !parsed.IsPrivateRegistry() {
		t.Error("open-vsx.org should not be the default registry")
	}
	built := purl.BuildPURLString("vscode-extension", "Muhammad-Sammy.CSharp", "2.15.30", "https://open-vsx.org")
	if built != "pkg:vscode-extension/muhammad-sammy/csharp@2.15.30?repository_url=https:%2F%2Fopen-vsx.org" {
		t.Errorf("BuildPURLString() = %q", built)
	}
}

func TestVSCodeExtensionRequiresPublisher(t *testing.T) {
	for _, name := range []string{"python", ".python", "ms-python.", "ms-python/python"} {
		if got := purl.MakePURLString("vscode-extension", name, "1.0.0"); got != "" {
			t.Errorf("MakePURLString(%q) = %q, want empty", name, got)
		}
		if got := purl.MakePURL("vscode-extension", name, "1.0.0"); got != nil {
			t.Errorf("MakePURL(%q) = %q, want nil", name, got)
		}
	}
}
