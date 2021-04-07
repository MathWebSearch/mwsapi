package xmlutil

import "github.com/beevik/etree"

// CopySafe is like element.Copy() except that it ensures all namespaces are set on the copied element.
func CopySafe(element *etree.Element) (copy *etree.Element) {
	if element == nil {
		return
	}

	// make a copy!
	copy = element.Copy()

	// ensure that the default namespace uri is set on clone!
	if elementNS, ok := GetDefaultNamespaceURI(element); ok {
		SetDefaultNamespaceURI(copy, elementNS)
	}

	// set all the dangling namespaces!
	for _, ns := range DanglingNamespaces(copy) {
		if uri, ok := GetNamespaceURI(element, ns); ok {
			SetNamespaceURI(copy, ns, uri)
		}
	}

	return
}
