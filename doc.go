// Package purl provides utilities for working with Package URLs (PURLs).
//
// It wraps github.com/package-url/packageurl-go with additional helpers
// for registry URL generation, type configuration, and version cleaning.
//
// # Parsing and Creating PURLs
//
//	// Parse a PURL string
//	p, err := purl.Parse("pkg:npm/%40babel/core@7.24.0")
//	fmt.Println(p.Type)      // npm
//	fmt.Println(p.Namespace) // @babel
//	fmt.Println(p.Name)      // core
//	fmt.Println(p.Version)   // 7.24.0
//	fmt.Println(p.FullName()) // @babel/core
//
//	// Create a PURL from components
//	p := purl.New("npm", "@babel", "core", "7.24.0", nil)
//	fmt.Println(p.String()) // pkg:npm/%40babel/core@7.24.0
//
// # Registry URLs
//
//	// Generate a registry URL from a PURL
//	p, _ := purl.Parse("pkg:npm/lodash@4.17.21")
//	url, _ := p.RegistryURLWithVersion()
//	fmt.Println(url) // https://www.npmjs.com/package/lodash/v/4.17.21
//
//	// Parse a registry URL back to a PURL
//	p, _ := purl.ParseRegistryURL("https://crates.io/crates/serde")
//	fmt.Println(p.String()) // pkg:cargo/serde
//
// # Type Configuration
//
// Type information comes from an embedded purl-types.json file.
//
//	purl.KnownTypes()           // []string of all known PURL types
//	purl.IsKnownType("npm")     // true
//	purl.DefaultRegistry("npm") // https://registry.npmjs.org
//
//	cfg := purl.TypeInfo("maven")
//	fmt.Println(cfg.NamespaceRequired()) // true
//
// # Private Registries
//
//	p, _ := purl.Parse("pkg:npm/lodash?repository_url=https://npm.example.com")
//	fmt.Println(p.IsPrivateRegistry()) // true
//	fmt.Println(p.RepositoryURL())     // https://npm.example.com
package purl
