package purl

import (
	_ "embed"
	"encoding/json"
	"sort"
	"sync"
)

//go:embed types.json
var typesJSON []byte

// TypeConfig contains configuration for a PURL type.
type TypeConfig struct {
	Description          string          `json:"description"`
	DefaultRegistry      *string         `json:"default_registry"`
	NamespaceRequirement string          `json:"namespace_requirement"`
	Examples             []string        `json:"examples"`
	RegistryConfig       *RegistryConfig `json:"registry_config"`
}

// RegistryConfig contains URL templates and patterns for registry URLs.
type RegistryConfig struct {
	BaseURL                         string              `json:"base_url"`
	ReverseRegex                    string              `json:"reverse_regex"`
	URITemplate                     string              `json:"uri_template"`
	URITemplateNoNamespace          string              `json:"uri_template_no_namespace"`
	URITemplateWithVersion          string              `json:"uri_template_with_version"`
	URITemplateWithVersionNoNS      string              `json:"uri_template_with_version_no_namespace"`
	Components                      RegistryComponents  `json:"components"`
}

// RegistryComponents describes which PURL components are used in registry URLs.
type RegistryComponents struct {
	Namespace         bool   `json:"namespace"`
	NamespaceRequired bool   `json:"namespace_required"`
	NamespacePrefix   string `json:"namespace_prefix"`
	VersionInURL      bool   `json:"version_in_url"`
	VersionPath       string `json:"version_path"`
	VersionPrefix     string `json:"version_prefix"`
	VersionSeparator  string `json:"version_separator"`
	DefaultVersion    string `json:"default_version"`
	TrailingSlash     bool   `json:"trailing_slash"`
	SpecialHandling   string `json:"special_handling"`
}

// NamespaceRequired returns true if the type requires a namespace.
func (t *TypeConfig) NamespaceRequired() bool {
	return t.NamespaceRequirement == "required"
}

// NamespaceProhibited returns true if the type prohibits namespaces.
func (t *TypeConfig) NamespaceProhibited() bool {
	return t.NamespaceRequirement == "prohibited"
}

type typesData struct {
	Version     string                `json:"version"`
	Description string                `json:"description"`
	Source      string                `json:"source"`
	LastUpdated string                `json:"last_updated"`
	Types       map[string]TypeConfig `json:"types"`
}

var (
	loadOnce   sync.Once
	loadedData *typesData
	loadErr    error
)

func loadTypes() (*typesData, error) {
	loadOnce.Do(func() {
		loadedData = &typesData{}
		loadErr = json.Unmarshal(typesJSON, loadedData)
	})
	return loadedData, loadErr
}

// TypeInfo returns configuration for a PURL type, or nil if unknown.
func TypeInfo(purlType string) *TypeConfig {
	data, err := loadTypes()
	if err != nil {
		return nil
	}
	cfg, ok := data.Types[purlType]
	if !ok {
		return nil
	}
	return &cfg
}

// KnownTypes returns a sorted list of all known PURL types.
func KnownTypes() []string {
	data, err := loadTypes()
	if err != nil {
		return nil
	}
	types := make([]string, 0, len(data.Types))
	for t := range data.Types {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// IsKnownType returns true if the PURL type is defined in types.json.
func IsKnownType(purlType string) bool {
	data, err := loadTypes()
	if err != nil {
		return false
	}
	_, ok := data.Types[purlType]
	return ok
}

// DefaultRegistry returns the default registry URL for a PURL type.
// Returns empty string if the type has no default registry.
func DefaultRegistry(purlType string) string {
	cfg := TypeInfo(purlType)
	if cfg == nil || cfg.DefaultRegistry == nil {
		return ""
	}
	return *cfg.DefaultRegistry
}

// TypesVersion returns the version of the types.json data.
func TypesVersion() string {
	data, err := loadTypes()
	if err != nil {
		return ""
	}
	return data.Version
}
