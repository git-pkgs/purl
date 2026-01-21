package purl

// FullName returns the package name combining namespace and name.
// If there's no namespace, returns just the name.
// The namespace and name are joined with "/".
func (p *PURL) FullName() string {
	if p.Namespace == "" {
		return p.Name
	}
	return p.Namespace + "/" + p.Name
}
