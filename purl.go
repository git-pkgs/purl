// Package purl provides utilities for working with Package URLs (PURLs).
//
// It wraps github.com/package-url/packageurl-go with additional helpers
// for ecosystem-specific name formatting, registry URL generation, and
// PURL construction from ecosystem-native package identifiers.
package purl

import (
	"errors"
	"sort"
	"strings"

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
	if p.Type == ecosystemGit {
		if p.Namespace == "" {
			return nil, errors.New("git purl requires a host and repository path")
		}
		if host, path, ok := strings.Cut(p.Namespace, "/"); ok {
			p.Namespace = host
			p.Name = path + "/" + p.Name
		}
	}
	lowercaseCaseInsensitive(&p)
	return &PURL{p}, nil
}

// lowercaseCaseInsensitive lowercases components the spec declares
// case-insensitive where packageurl-go's Normalize leaves them as given.
func lowercaseCaseInsensitive(p *packageurl.PackageURL) {
	if p.Type == purlTypeVSCodeExtension {
		p.Namespace = strings.ToLower(p.Namespace)
		p.Name = strings.ToLower(p.Name)
	}
}

// New creates a new PURL from components.
//
// The result is normalized with the same per-type rules Parse applies
// (packageurl-go's Normalize), so New and Parse produce the same string for
// equivalent inputs. Normalization errors are ignored; any adjustments
// Normalize applied before erroring are kept.
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
	_ = p.Normalize()
	lowercaseCaseInsensitive(p)
	return &PURL{*p}
}

// String returns the PURL as a string.
func (p *PURL) String() string {
	return p.ToString()
}

// ToString returns the PURL as a string.
func (p *PURL) ToString() string {
	components := p.PackageURL
	if components.Type == ecosystemGit {
		// Serialize repository path segments without escaping their separators.
		if i := strings.LastIndexByte(components.Name, '/'); i >= 0 {
			components.Namespace += "/" + components.Name[:i]
			components.Name = components.Name[i+1:]
		}
	}
	return components.String()
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
