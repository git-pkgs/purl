package purl

import (
	"regexp"
	"strings"

	packageurl "github.com/package-url/packageurl-go"
)

const (
	purlTypeConda = "conda"
	purlTypeNuGet = "nuget"
	purlTypeOCI   = "oci"
	purlTypePyPI  = "pypi"
)

var pypiSeparators = regexp.MustCompile(`[-_.]+`)

// ComparisonForm returns a copy of the PURL in a form where two spellings of
// the same package compare equal, following each ecosystem's own rules:
//
//   - pypi: name lowercased with runs of "-", "_" and "." folded to "-" (PEP 503)
//   - nuget: name and version lowercased
//   - composer: a leading "v" or "V" before a digit dropped from the version
//   - conda: a channel given as namespace or as a "channel/name" name moved
//     into the channel qualifier
//   - oci: a namespace folded into the name, so library/alpine is one name
//   - swift: registry identities (a scope without "." or "/") lowercased
//
// Parse and String are left as they are, since callers derive storage keys
// from them; use this form only to decide whether two PURLs name the same
// package or version.
func (p *PURL) ComparisonForm() *PURL {
	c := packageurl.PackageURL{
		Type:       p.Type,
		Namespace:  p.Namespace,
		Name:       p.Name,
		Version:    p.Version,
		Qualifiers: append(packageurl.Qualifiers(nil), p.Qualifiers...),
		Subpath:    p.Subpath,
	}
	switch c.Type {
	case purlTypePyPI:
		c.Name = pypiSeparators.ReplaceAllString(strings.ToLower(c.Name), "-")
	case purlTypeNuGet:
		c.Name, c.Version = strings.ToLower(c.Name), strings.ToLower(c.Version)
	case ecosystemComposer:
		// Packagist keeps a tag's "v"; Composer treats v6.1.7 and 6.1.7 as one version.
		if len(c.Version) > 1 && (c.Version[0] == 'v' || c.Version[0] == 'V') && isDigit(c.Version[1]) {
			c.Version = c.Version[1:]
		}
	case purlTypeConda:
		comparisonConda(&c)
	case purlTypeOCI:
		// The spec gives oci no namespace; a path-style repository is one name.
		if c.Namespace != "" {
			c.Namespace, c.Name = "", c.Namespace+"/"+c.Name
		}
	case ecosystemSwift:
		// Registry identities are case-insensitive (SE-0292). Source
		// coordinates such as github.com/apple are left alone.
		if c.Namespace != "" && !strings.ContainsAny(c.Namespace, "./") {
			c.Namespace, c.Name = strings.ToLower(c.Namespace), strings.ToLower(c.Name)
		}
	}
	return &PURL{c}
}

// comparisonConda moves the channel into the channel qualifier. A name
// written "channel/name" wins over a namespace, which wins over the qualifier.
func comparisonConda(c *packageurl.PackageURL) {
	channel := ""
	if c.Namespace != "" {
		channel, c.Namespace = c.Namespace, ""
	}
	if i := strings.LastIndex(c.Name, "/"); i >= 0 {
		channel, c.Name = c.Name[:i], c.Name[i+1:]
	}
	if channel == "" {
		return
	}
	for i, q := range c.Qualifiers {
		if q.Key == "channel" {
			c.Qualifiers[i].Value = channel
			return
		}
	}
	c.Qualifiers = append(c.Qualifiers, packageurl.Qualifier{Key: "channel", Value: channel})
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
