package purl

import (
	"strings"

	"github.com/git-pkgs/vers"
)

// CleanVersion extracts a version from a version constraint string.
// Uses the vers library to parse the constraint and extract the minimum bound.
// If parsing fails, returns the original string.
func CleanVersion(version, scheme string) string {
	if version == "" {
		return ""
	}

	r, err := vers.ParseNative(version, scheme)
	if err != nil || len(r.Intervals) == 0 {
		return version
	}

	// Return the minimum bound from the first interval
	if r.Intervals[0].Min != "" {
		return r.Intervals[0].Min
	}

	return version
}

// BuildPURLString builds a PURL string directly from ecosystem-native identifiers
// without creating intermediate PURL structs. This is the fast path for manifest
// parsing where we just need the string output.
func BuildPURLString(ecosystem, name, version, registryURL string) string {
	purlType := EcosystemToPURLType(ecosystem)
	cleanVersion := CleanVersion(version, purlType)
	namespace, pkgName := splitNamespace(ecosystem, name)

	if registryURL != "" && !IsNonDefaultRegistry(purlType, registryURL) {
		registryURL = ""
	}

	return buildPURLString(purlType, namespace, pkgName, cleanVersion, registryURL)
}

func buildPURLString(purlType, namespace, name, version, registryURL string) string {
	n := len("pkg:") + len(purlType) + escapedNamespaceLength(namespace) + 1 + escapedComponentLength(name)
	if version != "" {
		n += 1 + escapedComponentLength(version)
	}
	if registryURL != "" {
		n += len("?repository_url=") + escapedQualifierLength(registryURL)
	}

	var b strings.Builder
	b.Grow(n)

	b.WriteString("pkg:")
	b.WriteString(purlType)
	start := 0
	for i := 0; i <= len(namespace); i++ {
		if i == len(namespace) || namespace[i] == '/' {
			if i > start {
				b.WriteByte('/')
				writeComponentEscaped(&b, namespace[start:i])
			}
			start = i + 1
		}
	}
	b.WriteByte('/')
	writeComponentEscaped(&b, name)
	if version != "" {
		b.WriteByte('@')
		writeComponentEscaped(&b, version)
	}
	if registryURL != "" {
		b.WriteString("?repository_url=")
		writeQualifierEscaped(&b, registryURL)
	}

	return b.String()
}

func escapedNamespaceLength(namespace string) int {
	n := 0
	start := 0
	for i := 0; i <= len(namespace); i++ {
		if i == len(namespace) || namespace[i] == '/' {
			if i > start {
				n += 1 + escapedComponentLength(namespace[start:i])
			}
			start = i + 1
		}
	}
	return n
}

func escapedComponentLength(s string) int {
	n := len(s)
	for i := 0; i < len(s); i++ {
		if !isComponentSafe(s[i]) {
			n += 2 //nolint:mnd
		}
	}
	return n
}

func escapedQualifierLength(s string) int {
	n := len(s)
	for i := 0; i < len(s); i++ {
		if !isQualifierValueSafe(s[i]) {
			n += 2 //nolint:mnd
		}
	}
	return n
}

// splitNamespace extracts namespace and package name from an ecosystem-native
// package identifier.
func splitNamespace(ecosystem, name string) (namespace, pkgName string) {
	pkgName = name
	normalized := NormalizeEcosystem(ecosystem)

	if ns, ok := defaultNamespaces[normalized]; ok {
		namespace = ns
	}

	switch normalized {
	case ecosystemNPM:
		if strings.HasPrefix(name, "@") {
			if i := strings.IndexByte(name, '/'); i >= 0 {
				namespace = name[:i]
				pkgName = name[i+1:]
			}
		}
	case ecosystemGolang:
		if i := strings.LastIndex(name, "/"); i > 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case ecosystemMaven:
		if i := strings.IndexByte(name, ':'); i >= 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case ecosystemPackagist, ecosystemComposer:
		if i := strings.IndexByte(name, '/'); i >= 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case ecosystemGitHubActions:
		if i := strings.IndexByte(name, '/'); i >= 0 {
			namespace = name[:i]
			rest := name[i+1:]
			if j := strings.IndexByte(rest, '/'); j >= 0 {
				pkgName = rest[:j]
			} else {
				pkgName = rest
			}
		}
	}
	return
}

// writeComponentEscaped writes s to b, percent-encoding characters that are not safe
// in PURL path components (namespace, name, version).
func writeComponentEscaped(b *strings.Builder, s string) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isComponentSafe(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(hexDigit(c >> 4))   //nolint:mnd
			b.WriteByte(hexDigit(c & 0x0f)) //nolint:mnd
		}
	}
}

// isComponentSafe returns true for characters that can appear unencoded in
// PURL namespace/name/version segments. Matches the fork's isPurlSafe.
func isComponentSafe(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~' ||
		c == '!' || c == '$' || c == '&' || c == '\'' ||
		c == '(' || c == ')' || c == '*' ||
		c == ',' || c == ';' || c == '=' || c == ':'
}

// writeQualifierEscaped writes s to b, percent-encoding characters that are not safe
// in PURL qualifier values.
func writeQualifierEscaped(b *strings.Builder, s string) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isQualifierValueSafe(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(hexDigit(c >> 4))   //nolint:mnd
			b.WriteByte(hexDigit(c & 0x0f)) //nolint:mnd
		}
	}
}

func isQualifierValueSafe(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~' || c == ':'
}

func hexDigit(b byte) byte {
	if b < 10 { //nolint:mnd
		return '0' + b
	}
	return 'A' + b - 10 //nolint:mnd
}
