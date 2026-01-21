package purl

import (
	"testing"
)

func TestFullName(t *testing.T) {
	tests := []struct {
		purl string
		want string
	}{
		// Simple packages without namespace
		{"pkg:cargo/serde", "serde"},
		{"pkg:npm/lodash", "lodash"},
		{"pkg:pypi/requests", "requests"},
		{"pkg:gem/rails", "rails"},

		// npm scoped packages
		{"pkg:npm/%40babel/core@7.24.0", "@babel/core"},
		{"pkg:npm/%40types/node@18.0.0", "@types/node"},

		// Maven with groupId (uses : separator)
		{"pkg:maven/org.apache.commons/commons-lang3@3.12.0", "org.apache.commons:commons-lang3"},
		{"pkg:maven/junit/junit@4.13.2", "junit:junit"},

		// Go modules
		{"pkg:golang/github.com/gorilla/mux@v1.8.0", "github.com/gorilla/mux"},
		{"pkg:golang/google.golang.org/grpc@v1.50.0", "google.golang.org/grpc"},

		// Terraform modules
		{"pkg:terraform/hashicorp/consul/aws@0.11.0", "hashicorp/consul/aws"},

		// Composer
		{"pkg:composer/symfony/console@6.1.7", "symfony/console"},

		// Hex (no namespace)
		{"pkg:hex/phoenix@1.7.0", "phoenix"},

		// Elm
		{"pkg:elm/elm/http@2.0.0", "elm/http"},

		// Clojars
		{"pkg:clojars/org.clojure/clojure@1.11.1", "org.clojure/clojure"},
	}

	for _, tt := range tests {
		t.Run(tt.purl, func(t *testing.T) {
			p, err := Parse(tt.purl)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.purl, err)
			}
			if got := p.FullName(); got != tt.want {
				t.Errorf("FullName() = %q, want %q", got, tt.want)
			}
		})
	}
}
