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

	data, err := loadTypes()
	if err != nil {
		return false
	}
	cfg, ok := data.Types[purlType]
	if !ok || cfg.DefaultRegistry == nil {
		return false
	}

	defaultURL := *cfg.DefaultRegistry
	if defaultURL == "" {
		return false
	}

	defaultHost, ok := data.defaultHosts[defaultURL]
	if !ok {
		// TypeInfo exposes the default URL through a shared pointer.
		defaultHost = extractHost(defaultURL)
	}
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

// extractHost extracts the hostname from a URL string using net/url.Parse.
func extractHost(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
