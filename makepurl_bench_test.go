package purl

import (
	"fmt"
	"testing"
)

var benchmarkRegistryDefault bool

var resolvedURLCases = []struct {
	name     string
	registry string
	want     string
}{
	{"none", "", "pkg:npm/%40scope/package@1.2.3"},
	{"default", "https://registry.npmjs.org/@scope/package/-/package-1.2.3.tgz", "pkg:npm/%40scope/package@1.2.3"},
	{"subdomain", "https://cdn.registry.npmjs.org/@scope/package/-/package-1.2.3.tgz", "pkg:npm/%40scope/package@1.2.3"},
	{"private", "https://npm.example.invalid/@scope/package/-/package-1.2.3.tgz", "pkg:npm/%40scope/package@1.2.3?repository_url=https:%2F%2Fnpm.example.invalid%2F%40scope%2Fpackage%2F-%2Fpackage-1.2.3.tgz"},
}

func TestBuildPURLResolvedURLs(t *testing.T) {
	for _, tc := range resolvedURLCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildPURLString("npm", "@scope/package", "1.2.3", tc.registry); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func BenchmarkBuildPURLResolvedURL(b *testing.B) {
	for _, tc := range resolvedURLCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkVersionPURLString = BuildPURLString("npm", "@scope/package", "1.2.3", tc.registry)
			}
			if benchmarkVersionPURLString != tc.want {
				b.Fatalf("got %q, want %q", benchmarkVersionPURLString, tc.want)
			}
		})
	}
}

func BenchmarkBuildPURLResolvedURLMixed(b *testing.B) {
	const count = 1000
	type input struct{ name, url, want string }
	inputs := make([]input, count)
	for i := range inputs {
		name := fmt.Sprintf("@scope/package-%d", i)
		inputs[i] = input{
			name: name,
			url:  "https://registry.npmjs.org/" + name + "/-/package-1.2.3.tgz",
			want: fmt.Sprintf("pkg:npm/%%40scope/package-%d@1.2.3", i),
		}
		if got := BuildPURLString("npm", name, "1.2.3", inputs[i].url); got != inputs[i].want {
			b.Fatalf("got %q, want %q", got, inputs[i].want)
		}
	}
	i := 0
	b.ReportAllocs()
	for b.Loop() {
		in := inputs[i%count]
		benchmarkVersionPURLString = BuildPURLString("npm", in.name, "1.2.3", in.url)
		i++
	}
}

func BenchmarkDefaultRegistryResolvedURL(b *testing.B) {
	for _, tc := range resolvedURLCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkRegistryDefault = IsDefaultRegistry("npm", tc.registry)
			}
		})
	}
}

func BenchmarkCleanVersionNPM(b *testing.B) {
	for _, version := range []string{"1.2.3", "1.2.3-beta.1", "^1.2.3"} {
		b.Run(version, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkVersionPURLString = CleanVersion(version, "npm")
			}
		})
	}
}
