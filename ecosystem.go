package purl

import (
	"strings"
)

const (
	ecosystemAlpine        = "alpine"
	ecosystemArch          = "arch"
	ecosystemComposer      = "composer"
	ecosystemGitHubActions = "github-actions"
	ecosystemGolang        = "golang"
	ecosystemMaven         = "maven"
	ecosystemNPM           = "npm"
	ecosystemPackagist     = "packagist"
	ecosystemRubyGems      = "rubygems"
	ecosystemSwift         = "swift"
	purlTypeGem            = "gem"
	purlTypeGitHubActions  = "githubactions"
)

// purlTypeForEcosystem maps ecosystem names to PURL types.
// Most ecosystems use their name as the PURL type, but some differ.
var purlTypeForEcosystem = map[string]string{
	ecosystemAlpine:        "apk",
	ecosystemArch:          "alpm",
	ecosystemRubyGems:      purlTypeGem,
	ecosystemPackagist:     ecosystemComposer,
	ecosystemGitHubActions: purlTypeGitHubActions,
}

// ecosystemAliases maps alternate names to canonical ecosystem names.
var ecosystemAliases = map[string]string{
	"go":              ecosystemGolang,
	purlTypeGem:       ecosystemRubyGems,
	ecosystemComposer: ecosystemPackagist,
}

// osvEcosystemNames maps PURL types to OSV ecosystem names.
var osvEcosystemNames = map[string]string{
	purlTypeGem:           "RubyGems",
	ecosystemNPM:          ecosystemNPM,
	"pypi":                "PyPI",
	"cargo":               "crates.io",
	"conan":               "ConanCenter",
	"cran":                "CRAN",
	ecosystemGolang:       "Go",
	"hackage":             "Hackage",
	ecosystemMaven:        "Maven",
	"julia":               "Julia",
	"nuget":               "NuGet",
	"opam":                "opam",
	ecosystemComposer:     "Packagist",
	"hex":                 "Hex",
	"pub":                 "Pub",
	"swift":               "SwiftURL",
	purlTypeGitHubActions: "GitHub Actions",
}

// depsdevSystemNames maps PURL types to deps.dev system names.
var depsdevSystemNames = map[string]string{
	ecosystemNPM:    "NPM",
	purlTypeGem:     "RUBYGEMS",
	"pypi":          "PYPI",
	"cargo":         "CARGO",
	ecosystemGolang: "GO",
	ecosystemMaven:  "MAVEN",
	"nuget":         "NUGET",
}

// defaultNamespaces defines default namespaces for certain ecosystems.
var defaultNamespaces = map[string]string{
	ecosystemAlpine: ecosystemAlpine,
	ecosystemArch:   ecosystemArch,
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

// PURLTypeToOSV converts a PURL type to the OSV ecosystem name and reports
// whether the type is one OSV recognises. Unlike EcosystemToOSV, it takes the
// PURL type directly (skipping the ecosystem-name normalisation when the
// caller already has a parsed PURL) and returns ok=false on a miss rather
// than passing the input through, so callers that emit OSV records can fall
// back to a GIT range instead of writing an ecosystem the OSV schema will
// reject.
func PURLTypeToOSV(purlType string) (string, bool) {
	osv, ok := osvEcosystemNames[purlType]
	return osv, ok
}

// PURLTypeToDepsdev converts a PURL type to the deps.dev system name.
// Returns empty string if the type is not supported by deps.dev.
func PURLTypeToDepsdev(purlType string) string {
	if system, ok := depsdevSystemNames[purlType]; ok {
		return system
	}
	return ""
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
//   - swift: host/owner/package -> namespace="host/owner", name="package"
//
// Swift registry identities do not contain source repository coordinates and
// return nil because the Swift PURL type cannot represent them.
func MakePURL(ecosystem, name, version string) *PURL {
	purlType := EcosystemToPURLType(ecosystem)
	namespace, pkgName, ok := splitNamespace(ecosystem, name)
	if !ok {
		return nil
	}

	return New(purlType, namespace, pkgName, version, nil)
}

// MakePURLString is like MakePURL but returns the PURL as a string. It returns
// an empty string when the package identifier cannot be represented as a PURL.
func MakePURLString(ecosystem, name, version string) string {
	purlType := EcosystemToPURLType(ecosystem)
	namespace, pkgName, ok := splitNamespace(ecosystem, name)
	if !ok {
		return ""
	}
	return buildPURLString(purlType, namespace, pkgName, version, "")
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
