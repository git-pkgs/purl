// Package purl provides utilities for working with Package URLs (PURLs).
//
// It wraps github.com/package-url/packageurl-go with additional helpers
// for ecosystem-specific name formatting, registry URL generation, and
// PURL construction from ecosystem-native package identifiers.
package purl

import (
	"sort"

	packageurl "github.com/package-url/packageurl-go"
)

// PURL wraps packageurl.PackageURL with additional helpers.
type PURL struct {
	packageurl.PackageURL
}

// Parse parses a Package URL string into a PURL.
func Parse(s string) (*PURL, error) {
	p, err := packageurl.FromString(s)
	if err != nil {
		return nil, err
	}
	return &PURL{p}, nil
}

// New creates a new PURL from components.
func New(purlType, namespace, name, version string, qualifiers map[string]string) *PURL {
	var q packageurl.Qualifiers
	if len(qualifiers) > 0 {
		// Sort keys for deterministic output
		keys := make([]string, 0, len(qualifiers))
		for k := range qualifiers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			q = append(q, packageurl.Qualifier{Key: k, Value: qualifiers[k]})
		}
	}
	p := packageurl.NewPackageURL(purlType, namespace, name, version, q, "")
	return &PURL{*p}
}

// String returns the PURL as a string.
func (p *PURL) String() string {
	return p.PackageURL.String()
}

// RepositoryURL returns the repository_url qualifier value, if present.
func (p *PURL) RepositoryURL() string {
	return p.Qualifiers.Map()["repository_url"]
}

// IsPrivateRegistry returns true if the PURL has a non-default repository_url.
func (p *PURL) IsPrivateRegistry() bool {
	repoURL := p.RepositoryURL()
	if repoURL == "" {
		return false
	}
	return IsNonDefaultRegistry(p.Type, repoURL)
}

// Qualifier returns the value of a qualifier, or empty string if not present.
func (p *PURL) Qualifier(key string) string {
	return p.Qualifiers.Map()[key]
}

// WithVersion returns a copy of the PURL with a different version.
func (p *PURL) WithVersion(version string) *PURL {
	return &PURL{
		PackageURL: packageurl.PackageURL{
			Type:       p.Type,
			Namespace:  p.Namespace,
			Name:       p.Name,
			Version:    version,
			Qualifiers: p.Qualifiers,
			Subpath:    p.Subpath,
		},
	}
}

// WithoutVersion returns a copy of the PURL without a version.
func (p *PURL) WithoutVersion() *PURL {
	return p.WithVersion("")
}

// WithQualifier returns a copy of the PURL with the qualifier set.
// If the qualifier already exists, it is replaced.
func (p *PURL) WithQualifier(key, value string) *PURL {
	q := make(packageurl.Qualifiers, 0, len(p.Qualifiers)+1)
	replaced := false
	for _, qual := range p.Qualifiers {
		if qual.Key == key {
			q = append(q, packageurl.Qualifier{Key: key, Value: value})
			replaced = true
		} else {
			q = append(q, qual)
		}
	}
	if !replaced {
		q = append(q, packageurl.Qualifier{Key: key, Value: value})
	}
	return &PURL{
		PackageURL: packageurl.PackageURL{
			Type:       p.Type,
			Namespace:  p.Namespace,
			Name:       p.Name,
			Version:    p.Version,
			Qualifiers: q,
			Subpath:    p.Subpath,
		},
	}
}
