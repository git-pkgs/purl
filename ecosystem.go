package purl

import (
	"strings"
)

// purlTypeForEcosystem maps ecosystem names to PURL types.
// Most ecosystems use their name as the PURL type, but some differ.
var purlTypeForEcosystem = map[string]string{
	"alpine":    "apk",
	"arch":      "alpm",
	"rubygems":  "gem",
	"packagist": "composer",
}

// ecosystemAliases maps alternate names to canonical ecosystem names.
var ecosystemAliases = map[string]string{
	"go":        "golang",
	"gem":       "rubygems",
	"composer":  "packagist",
}

// osvEcosystemNames maps PURL types to OSV ecosystem names.
var osvEcosystemNames = map[string]string{
	"gem":       "RubyGems",
	"npm":       "npm",
	"pypi":      "PyPI",
	"cargo":     "crates.io",
	"golang":    "Go",
	"maven":     "Maven",
	"nuget":     "NuGet",
	"composer":  "Packagist",
	"hex":       "Hex",
	"pub":       "Pub",
	"cocoapods": "CocoaPods",
}

// defaultNamespaces defines default namespaces for certain ecosystems.
var defaultNamespaces = map[string]string{
	"alpine": "alpine",
	"arch":   "arch",
}

// NormalizeEcosystem returns the canonical ecosystem name.
// Handles aliases like "go" -> "golang", "gem" -> "rubygems".
func NormalizeEcosystem(ecosystem string) string {
	lower := strings.ToLower(ecosystem)
	if canonical, ok := ecosystemAliases[lower]; ok {
		return canonical
	}
	return lower
}

// EcosystemToPURLType converts an ecosystem name to the corresponding PURL type.
// Returns the input unchanged if no mapping exists.
func EcosystemToPURLType(ecosystem string) string {
	normalized := NormalizeEcosystem(ecosystem)
	if t, ok := purlTypeForEcosystem[normalized]; ok {
		return t
	}
	return normalized
}

// PURLTypeToEcosystem converts a PURL type back to an ecosystem name.
// This is the inverse of EcosystemToPURLType.
func PURLTypeToEcosystem(purlType string) string {
	// Reverse lookup
	for eco, pt := range purlTypeForEcosystem {
		if pt == purlType {
			return eco
		}
	}
	return purlType
}

// EcosystemToOSV converts an ecosystem name to the OSV ecosystem name.
// OSV uses specific capitalization and naming conventions.
func EcosystemToOSV(ecosystem string) string {
	purlType := EcosystemToPURLType(ecosystem)
	if osv, ok := osvEcosystemNames[purlType]; ok {
		return osv
	}
	return ecosystem
}

// MakePURL constructs a PURL from ecosystem-native package identifiers.
//
// It handles namespace extraction for ecosystems:
//   - npm: @scope/pkg -> namespace="@scope", name="pkg"
//   - maven: group:artifact -> namespace="group", name="artifact"
//   - golang: github.com/foo/bar -> namespace="github.com/foo", name="bar"
//   - composer: vendor/package -> namespace="vendor", name="package"
//   - alpine: pkg -> namespace="alpine", name="pkg"
//   - arch: pkg -> namespace="arch", name="pkg"
func MakePURL(ecosystem, name, version string) *PURL {
	purlType := EcosystemToPURLType(ecosystem)
	namespace := ""
	pkgName := name

	// Handle default namespaces
	if ns, ok := defaultNamespaces[NormalizeEcosystem(ecosystem)]; ok {
		namespace = ns
	}

	// Extract namespace from name based on ecosystem conventions
	switch NormalizeEcosystem(ecosystem) {
	case "npm":
		if strings.HasPrefix(name, "@") {
			parts := strings.SplitN(name, "/", 2)
			if len(parts) == 2 {
				namespace = parts[0] // Keep the @ for packageurl-go
				pkgName = parts[1]
			}
		}
	case "golang":
		if idx := strings.LastIndex(name, "/"); idx > 0 {
			namespace = name[:idx]
			pkgName = name[idx+1:]
		}
	case "maven":
		if strings.Contains(name, ":") {
			parts := strings.SplitN(name, ":", 2)
			namespace = parts[0]
			pkgName = parts[1]
		}
	case "packagist", "composer":
		if strings.Contains(name, "/") {
			parts := strings.SplitN(name, "/", 2)
			namespace = parts[0]
			pkgName = parts[1]
		}
	}

	return New(purlType, namespace, pkgName, version, nil)
}

// MakePURLString is like MakePURL but returns the PURL as a string.
func MakePURLString(ecosystem, name, version string) string {
	return MakePURL(ecosystem, name, version).String()
}

// SupportedEcosystems returns a list of all supported ecosystem names.
// This includes both PURL types and common aliases.
func SupportedEcosystems() []string {
	seen := make(map[string]bool)
	var result []string

	// Add all known PURL types
	for _, t := range KnownTypes() {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}

	// Add ecosystem aliases
	for alias := range ecosystemAliases {
		if !seen[alias] {
			seen[alias] = true
			result = append(result, alias)
		}
	}

	// Add ecosystems that map to different PURL types
	for eco := range purlTypeForEcosystem {
		if !seen[eco] {
			seen[eco] = true
			result = append(result, eco)
		}
	}

	return result
}

// IsValidEcosystem returns true if the ecosystem is recognized.
func IsValidEcosystem(ecosystem string) bool {
	normalized := NormalizeEcosystem(ecosystem)
	purlType := EcosystemToPURLType(normalized)
	return IsKnownType(purlType)
}
