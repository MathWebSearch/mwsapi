package xmlutil

import (
	"sort"

	"github.com/beevik/etree"
)

// DanglingNamespacesMap recursively scans element for used namespaces in either Elements or Attributes.
// It returns a sorted list of namespaces that are used, but not explicitly declared.
// DanglingNamespacesMap respects the "xml" and "xmlns" predefined namespaces.
func DanglingNamespaces(element *etree.Element) []string {
	// find the referenced elements!
	referencedMap := DanglingNamespacesMap(element, []string{})
	delete(referencedMap, "xml")
	delete(referencedMap, "xmlns")

	// turn them into an array!
	var i int
	referenced := make([]string, len(referencedMap))
	for ns := range referencedMap {
		referenced[i] = ns
		i++
	}

	// sort and return them!
	sort.Strings(referenced)

	return referenced
}

// DanglingNamespacesMap recursively scans element for used namespaces in either Elements or Attributes.
// It returns a map of namespaces that are used, but not explicitly declared.
// DanglingNamespacesMap does not respect predefined namespaces "xml" and "xmlns".
//
// Skip is a list of namespace URIs that can be ignored.
// Note that this function will call
func DanglingNamespacesMap(element *etree.Element, skip []string) map[string]struct{} {
	// a nil element does not have any namespacesq
	if element == nil {
		return nil
	}

	// a list of used namespaces
	used := make(map[string]struct{})

	// check if the element has a namespace
	if element.Space != "" {
		used[element.Space] = struct{}{}
	}

	// check if the attributes have a namespace
	for _, a := range element.Attr {
		// used the declaration
		if a.Space != "" {
			used[a.Key] = struct{}{} // make the value as used!
		}

		// mark the value as declared
		if a.Space == "xmlns" && a.Key != "" {
			skip = append(skip, a.Key)
		}
	}

	// all the keys we used are ok!
	for _, key := range skip {
		delete(used, key)
	}

	// add the keys in any of the children
	for _, child := range element.ChildElements() {
		for ns := range DanglingNamespacesMap(child, skip) {
			used[ns] = struct{}{}
		}
	}

	return used
}

// GetNamespaceURI finds a local namespace uri at an element.
// Adatped from etree.findLocalNamespaceURI() internal.
func GetNamespaceURI(element *etree.Element, space string) (value string, ok bool) {
	for _, a := range element.Attr {
		if a.Space == "xmlns" && a.Key == space {
			return a.Value, true
		}
	}

	parent := element.Parent()
	if parent == nil {
		return "", false
	}

	return GetNamespaceURI(parent, space)
}

// SetNamespaceURI sets the namespace space to value on element and all its' children.
// It modifies element in place, but does not modify parent values (if any).
func SetNamespaceURI(element *etree.Element, space, value string) {
	for i, a := range element.Attr {
		if a.Space == "xmlns" && a.Key == space {
			element.Attr[i].Value = value
			return
		}
	}

	element.CreateAttr("xmlns:"+space, value)
}

// GetDefaultNamespaceURI finds the default namespace of an element.
// Adapted from etree.findDefaultNamespaceURI() internal
func GetDefaultNamespaceURI(element *etree.Element) (value string, ok bool) {
	for _, a := range element.Attr {
		if a.Space == "" && a.Key == "xmlns" {
			return a.Value, true
		}
	}

	parent := element.Parent()
	if parent == nil {
		return "", false
	}

	return GetDefaultNamespaceURI(parent)
}

// SetDefaultNamespaceURI sets the default namespace of element and all its' children.
// It modifies element in place, but does not modify parent values (if any).
func SetDefaultNamespaceURI(element *etree.Element, value string) {
	for i, a := range element.Attr {
		if a.Space == "" && a.Key == "xmlns" {
			element.Attr[i].Value = value
			return
		}
	}

	element.CreateAttr("xmlns", value)
}
