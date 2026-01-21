package purl

import (
	"net/url"
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

// extractHost parses a URL and returns its host.
func extractHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
