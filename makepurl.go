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

	needsQualifier := registryURL != "" && IsNonDefaultRegistry(purlType, registryURL)

	// Estimate capacity
	n := 4 + len(purlType) + 1 + len(pkgName) // "pkg:" + type + "/" + name
	if namespace != "" {
		n += 1 + len(namespace) // "/" + namespace
	}
	if cleanVersion != "" {
		n += 1 + len(cleanVersion) // "@" + version
	}
	if needsQualifier {
		n += 16 + len(registryURL) // "?repository_url=" + value
	}

	var b strings.Builder
	b.Grow(n)

	b.WriteString("pkg:")
	b.WriteString(purlType)
	if namespace != "" {
		// Write namespace segments, escaping each one
		for namespace != "" {
			b.WriteByte('/')
			seg := namespace
			if i := strings.IndexByte(namespace, '/'); i >= 0 {
				seg = namespace[:i]
				namespace = namespace[i+1:]
			} else {
				namespace = ""
			}
			writeComponentEscaped(&b, seg)
		}
	}
	b.WriteByte('/')
	writeComponentEscaped(&b, pkgName)
	if cleanVersion != "" {
		b.WriteByte('@')
		writeComponentEscaped(&b, cleanVersion)
	}
	if needsQualifier {
		b.WriteString("?repository_url=")
		writeQualifierEscaped(&b, registryURL)
	}

	return b.String()
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
	case "npm":
		if strings.HasPrefix(name, "@") {
			if i := strings.IndexByte(name, '/'); i >= 0 {
				namespace = name[:i]
				pkgName = name[i+1:]
			}
		}
	case "golang":
		if i := strings.LastIndex(name, "/"); i > 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case "maven":
		if i := strings.IndexByte(name, ':'); i >= 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case "packagist", "composer":
		if i := strings.IndexByte(name, '/'); i >= 0 {
			namespace = name[:i]
			pkgName = name[i+1:]
		}
	case "github-actions":
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
			b.WriteByte(hexDigit(c >> 4))
			b.WriteByte(hexDigit(c & 0x0f))
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
			b.WriteByte(hexDigit(c >> 4))
			b.WriteByte(hexDigit(c & 0x0f))
		}
	}
}

func isQualifierValueSafe(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~' || c == ':'
}

func hexDigit(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'A' + b - 10
}
