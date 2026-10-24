package purl_test

import (
	"testing"

	"github.com/git-pkgs/purl"
)

func TestDockerImagePURLs(t *testing.T) {
	tests := []struct {
		name, version, want string
	}{
		{"postgres", "16", "pkg:docker/postgres@16"},
		{"example/image", "latest", "pkg:docker/example/image@latest"},
		{"ghcr.io/example/image", "2", "pkg:docker/ghcr.io/example/image@2"},
		{"registry.example:5000/team/image", "1", "pkg:docker/registry.example:5000/team/image@1"},
		{"ghcr.io/example/image", "sha256:abc123", "pkg:docker/ghcr.io/example/image@sha256:abc123"},
		{"ghcr.io/example/image", "", "pkg:docker/ghcr.io/example/image"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"@"+tt.version, func(t *testing.T) {
			if got := purl.BuildPURLString("docker", tt.name, tt.version, ""); got != tt.want {
				t.Errorf("BuildPURLString = %q, want %q", got, tt.want)
			}
			if got := purl.MakePURLString("docker", tt.name, tt.version); got != tt.want {
				t.Errorf("MakePURLString = %q, want %q", got, tt.want)
			}
			got := purl.MakePURL("docker", tt.name, tt.version)
			if got == nil {
				t.Fatal("MakePURL returned nil")
			}
			if got.String() != tt.want {
				t.Errorf("MakePURL = %q, want %q", got.String(), tt.want)
			}
		})
	}
}
