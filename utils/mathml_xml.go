package utils

import "github.com/antchfx/xmlquery"

// this file contains utility functions
// used by the mathml code

// findElementByLocalName finds an element by local name
func findElementByLocalName(node *xmlquery.Node, name string) *xmlquery.Node {
	if node == nil {
		return nil
	}

	// if we found the semantics node, return it
	if node.Type == xmlquery.ElementNode && node.Data == name {
		return node
	}

	// iterate over the children
	child := node.FirstChild
	var result *xmlquery.Node
	for child != nil {
		if result = findElementByLocalName(child, name); result != nil {
			return result
		}
		child = child.NextSibling
	}

	return nil
}

// finds the children of a <semantics> node that represent presentation and content mathml accordingly
func findPresentationAndContent(node *xmlquery.Node) (semantics *xmlquery.Node, annotation *xmlquery.Node) {
	child := node.FirstChild
	for child != nil && (annotation == nil || semantics == nil) {
		if child.Type == xmlquery.ElementNode {
			if child.Data == "annotation-xml" && annotation == nil {
				for _, attr := range child.Attr {
					if attr.Name.Local == "encoding" && attr.Value == "MathML-Content" {
						annotation = child
					}
				}
			} else if child.Data != "annotation" && semantics == nil {
				semantics = child
			}
		}
		child = child.NextSibling
	}
	return
}

// firstChildElement returns the first child of a node that is an element
// or nil if there are none
func firstChildElement(node *xmlquery.Node) *xmlquery.Node {
	child := node.FirstChild
	for child != nil {
		if child.Type == xmlquery.ElementNode {
			return child
		}
		child = child.NextSibling
	}
	return nil
}

// findElementByID finds an element by id
func findElementByID(node *xmlquery.Node, id string) *xmlquery.Node {
	if node == nil {
		return nil
	}

	// check the node itself
	if node.Type == xmlquery.ElementNode {
		for _, attr := range node.Attr {
			if attr.Name.Local == "id" && attr.Value == id {
				return node
			}
		}
	}

	// iterate over the children
	child := node.FirstChild
	var result *xmlquery.Node
	for child != nil {
		if result = findElementByID(child, id); result != nil {
			return result
		}
		child = child.NextSibling
	}

	return nil
}
