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
```

## Type Configuration

Type information comes from an embedded copy of [purl-types.json](https://github.com/andrew/purl/blob/main/purl-types.json) which includes default registries, namespace requirements, and URI templates for registry URL generation.

```go
purl.KnownTypes()           // []string of all known PURL types
purl.IsKnownType("npm")     // true
purl.DefaultRegistry("npm") // https://registry.npmjs.org
```

## License

MIT
