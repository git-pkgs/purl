package purl

// FullName returns the package name combining namespace and name.
// If there's no namespace, returns just the name.
// Maven uses ":" as separator (groupId:artifactId), all others use "/".
func (p *PURL) FullName() string {
	if p.Namespace == "" {
		return p.Name
	}
	if p.Type == "maven" {
		return p.Namespace + ":" + p.Name
	}
	return p.Namespace + "/" + p.Name
}
