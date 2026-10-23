# purl

Go library for working with Package URLs (PURLs). Wraps [packageurl-go](https://github.com/package-url/packageurl-go) with additional helpers for registry URL generation, type configuration, and version cleaning.

## Installation

```
go get github.com/git-pkgs/purl
```

## Usage

```go
import "github.com/git-pkgs/purl"

// Parse a PURL
p, _ := purl.Parse("pkg:npm/%40babel/core@7.24.0")
fmt.Println(p.FullName())  // @babel/core
fmt.Println(p.Type)        // npm
fmt.Println(p.Version)     // 7.24.0

// Create a PURL
p := purl.New("npm", "@babel", "core", "7.24.0", nil)
fmt.Println(p.String())  // pkg:npm/%40babel/core@7.24.0

// Generate registry URL
url, _ := p.RegistryURL()
fmt.Println(url)  // https://www.npmjs.com/package/@babel/core

// Parse registry URL back to PURL
p, _ := purl.ParseRegistryURL("https://crates.io/crates/serde")
fmt.Println(p.String())  // pkg:cargo/serde

// Check if using private registry
p, _ := purl.Parse("pkg:npm/lodash?repository_url=https://npm.example.com")
fmt.Println(p.IsPrivateRegistry())  // true

// Get type configuration
cfg := purl.TypeInfo("npm")
fmt.Println(*cfg.DefaultRegistry)  // https://registry.npmjs.org

// Clean version constraints
version := purl.CleanVersion("^1.2.3", "npm")
fmt.Println(version)  // 1.2.3

// Compare two spellings of the same package
a, _ := purl.Parse("pkg:pypi/zope.interface@6.0")
b, _ := purl.Parse("pkg:pypi/Zope_Interface@6.0")
fmt.Println(a.ComparisonForm().String() == b.ComparisonForm().String())  // true
```

`ComparisonForm` applies ecosystem rules that `Parse` does not, so `String()`
output stays stable for callers that store it:

| Type | Rule |
|------|------|
| pypi | name lowercased, runs of `-`, `_`, `.` become `-` (PEP 503) |
| nuget | name and version lowercased |
| composer | leading `v`/`V` before a digit dropped from the version |
| conda | channel given as namespace or `channel/name` moved to the `channel` qualifier |
| oci | namespace folded into the name (`library/alpine`) |
| swift | registry identities (scope without `.` or `/`) lowercased |

## Type Configuration

Type information is embedded in `types.json`. Type names and namespace requirements are checked against the pinned [purl-spec](https://github.com/package-url/purl-spec) submodule. Default registries and URL templates are maintained locally.

The local types `bioconductor`, `clojars`, `deno`, and `elm` are retained. `bsd` follows the FreeBSD type proposed in [purl-spec#496](https://github.com/package-url/purl-spec/pull/496); the `freebsd` ecosystem maps to `pkg:bsd/freebsd/...`. The legacy `homebrew` type keeps its existing PURLs and registry behavior; use `brew` explicitly for the registered type. Recognition of a type does not imply full support for its normalization or validation rules.

```go
purl.KnownTypes()           // []string of all known PURL types
purl.IsKnownType("npm")     // true
purl.DefaultRegistry("npm") // https://registry.npmjs.org
```

## Development

Initialize the spec submodule before running tests:

```sh
git submodule update --init
go test ./...
```

Dependabot checks for spec updates weekly. The metadata test reports missing types, namespace differences, and unregistered types without a local allowance. Update `types.json` when those checks fail, preserving local registry settings.

## License

[MIT](LICENSE).
