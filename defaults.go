package purl

import (
	"strings"
)

// IsDefaultRegistry returns true if the registryURL matches the default registry for the type.
func IsDefaultRegistry(purlType, registryURL string) bool {
	if registryURL == "" {
		return true
	}

	cfg := TypeInfo(purlType)
	if cfg == nil || cfg.DefaultRegistry == nil {
		return false
	}

	defaultURL := *cfg.DefaultRegistry
	if defaultURL == "" {
		return false
	}

	// Compare hosts
	defaultHost := extractHost(defaultURL)
	givenHost := extractHost(registryURL)

	if defaultHost == "" || givenHost == "" {
		return false
	}

	return givenHost == defaultHost || strings.HasSuffix(givenHost, "."+defaultHost)
}

// IsNonDefaultRegistry returns true if the registryURL is not the default registry for the type.
func IsNonDefaultRegistry(purlType, registryURL string) bool {
	if registryURL == "" {
		return false
	}
	return !IsDefaultRegistry(purlType, registryURL)
}

// extractHost extracts the host from a URL string without using net/url.Parse.
func extractHost(rawURL string) string {
	s := rawURL
	// Strip scheme
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	// Strip userinfo
	if i := strings.Index(s, "@"); i >= 0 {
		s = s[i+1:]
	}
	// Strip path, query, fragment
	for _, sep := range []byte{'/', '?', '#'} {
		if i := strings.IndexByte(s, sep); i >= 0 {
			s = s[:i]
		}
	}
	return s
}
