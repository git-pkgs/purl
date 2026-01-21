package purl

import (
	"errors"
	"regexp"
	"strings"
	"sync"
)

// ErrNoRegistryConfig is returned when a PURL type has no registry configuration.
var ErrNoRegistryConfig = errors.New("no registry configuration for this type")

// ErrNoMatch is returned when a URL doesn't match the reverse regex.
var ErrNoMatch = errors.New("URL does not match any known registry pattern")

// regexCache caches compiled regular expressions by pattern.
var regexCache sync.Map

// RegistryURL returns the human-readable registry URL for the package.
// For example, pkg:npm/lodash returns "https://www.npmjs.com/package/lodash".
func (p *PURL) RegistryURL() (string, error) {
	cfg := TypeInfo(p.Type)
	if cfg == nil || cfg.RegistryConfig == nil {
		return "", ErrNoRegistryConfig
	}

	return expandTemplate(cfg.RegistryConfig, p.Namespace, p.Name, "")
}

// RegistryURLWithVersion returns the registry URL including version.
// Falls back to RegistryURL if version URLs aren't supported.
func (p *PURL) RegistryURLWithVersion() (string, error) {
	if p.Version == "" {
		return p.RegistryURL()
	}

	cfg := TypeInfo(p.Type)
	if cfg == nil || cfg.RegistryConfig == nil {
		return "", ErrNoRegistryConfig
	}

	return expandTemplate(cfg.RegistryConfig, p.Namespace, p.Name, p.Version)
}

// expandTemplate expands a URI template with the given components.
func expandTemplate(rc *RegistryConfig, namespace, name, version string) (string, error) {
	var template string

	hasNamespace := namespace != ""

	// Select the appropriate template
	if version != "" && rc.Components.VersionInURL {
		if hasNamespace && rc.URITemplateWithVersion != "" {
			template = rc.URITemplateWithVersion
		} else if !hasNamespace && rc.URITemplateWithVersionNoNS != "" {
			template = rc.URITemplateWithVersionNoNS
		} else if rc.URITemplateWithVersion != "" {
			template = rc.URITemplateWithVersion
		}
	}

	if template == "" {
		if hasNamespace {
			template = rc.URITemplate
		} else if rc.URITemplateNoNamespace != "" {
			template = rc.URITemplateNoNamespace
		} else {
			template = rc.URITemplate
		}
	}

	if template == "" {
		return "", ErrNoRegistryConfig
	}

	// Handle namespace prefix (e.g., @ for npm)
	displayNamespace := namespace
	if rc.Components.NamespacePrefix != "" && namespace != "" {
		// Add prefix if not already present
		if !strings.HasPrefix(namespace, rc.Components.NamespacePrefix) {
			displayNamespace = rc.Components.NamespacePrefix + namespace
		}
	}

	// Expand template variables
	result := template
	result = strings.ReplaceAll(result, "{namespace}", displayNamespace)
	result = strings.ReplaceAll(result, "{name}", name)
	result = strings.ReplaceAll(result, "{version}", version)

	return result, nil
}

// ParseRegistryURL attempts to parse a registry URL into a PURL.
// It tries all known types to find a match.
func ParseRegistryURL(url string) (*PURL, error) {
	for _, t := range KnownTypes() {
		p, err := ParseRegistryURLWithType(url, t)
		if err == nil {
			return p, nil
		}
	}
	return nil, ErrNoMatch
}

// ParseRegistryURLWithType parses a registry URL using a specific PURL type.
func ParseRegistryURLWithType(url, purlType string) (*PURL, error) {
	cfg := TypeInfo(purlType)
	if cfg == nil || cfg.RegistryConfig == nil || cfg.RegistryConfig.ReverseRegex == "" {
		return nil, ErrNoRegistryConfig
	}

	re, err := getOrCompileRegex(cfg.RegistryConfig.ReverseRegex)
	if err != nil {
		return nil, err
	}

	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return nil, ErrNoMatch
	}

	var namespace, name, version string

	// Parse matches based on component configuration
	if cfg.RegistryConfig.Components.Namespace {
		if cfg.RegistryConfig.Components.NamespaceRequired {
			// Namespace is required: matches[1]=namespace, matches[2]=name, matches[3]=version (if present)
			if len(matches) > 1 {
				namespace = matches[1]
			}
			if len(matches) > 2 {
				name = matches[2]
			}
			if len(matches) > 3 {
				version = matches[3]
			}
		} else {
			// Namespace is optional: matches[1]=namespace (maybe empty), matches[2]=name
			if len(matches) > 2 {
				namespace = matches[1]
				name = matches[2]
			}
			if len(matches) > 3 {
				version = matches[3]
			}
		}
	} else {
		// No namespace: matches[1]=name, matches[2]=version (if present)
		if len(matches) > 1 {
			name = matches[1]
		}
		if len(matches) > 2 {
			version = matches[2]
		}
	}

	if name == "" {
		return nil, ErrNoMatch
	}

	return New(purlType, namespace, name, version, nil), nil
}

// getOrCompileRegex returns a cached compiled regex or compiles and caches it.
func getOrCompileRegex(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexCache.Store(pattern, re)
	return re, nil
}
