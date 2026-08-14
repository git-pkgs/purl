package purl

import "testing"

var (
	benchmarkPackagePURLString string
	benchmarkVersionPURLString string
)

func BenchmarkMakePURLString(b *testing.B) {
	benchmarks := []struct {
		name      string
		ecosystem string
		pkg       string
		version   string
	}{
		{name: "npm", ecosystem: "npm", pkg: "lodash", version: "4.17.21"},
		{name: "npm_scoped", ecosystem: "npm", pkg: "@babel/core", version: "7.24.0"},
		{name: "golang", ecosystem: "golang", pkg: "github.com/foo/bar", version: "v1.0.0"},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkPackagePURLString = MakePURLString(
					benchmark.ecosystem,
					benchmark.pkg,
					"",
				)
				benchmarkVersionPURLString = MakePURLString(
					benchmark.ecosystem,
					benchmark.pkg,
					benchmark.version,
				)
			}
		})
	}
}
