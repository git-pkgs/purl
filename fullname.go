package purl

import "strings"

// FullName returns the package name combining namespace and name.
// If there's no namespace, returns just the name.
// Maven uses ":" and Swift registry identities use "." as the separator.
// Swift source coordinates and other types use "/".
func (p *PURL) FullName() string {
	if p.Namespace == "" {
		return p.Name
	}
	if p.Type == "maven" {
		return p.Namespace + ":" + p.Name
	}
	if p.Type == ecosystemSwift && !strings.ContainsAny(p.Namespace, "./") {
		return p.Namespace + "." + p.Name
	}
	return p.Namespace + "/" + p.Name
}
