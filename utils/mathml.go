package utils

import (
	"encoding/xml"
	"strings"

	"github.com/pkg/errors"

	"github.com/antchfx/xpath"

	"github.com/antchfx/xmlquery"
)

// MathML represents a MathML Element with semantics and annotation component
type MathML struct {
	root *xmlquery.Node

	semantics  *xmlquery.Node
	annotation *xmlquery.Node
}

// MathMLNamespace is the xml namespace of mathml
const MathMLNamespace = "http://www.w3.org/1998/Math/MathML"

// ParseMathML parses a new MathML Element
func ParseMathML(source string) (math *MathML, err error) {
	// wrap safe source in an xml header and something to add the <m> namespace
	reader := strings.NewReader(xml.Header + "<wrapper xmlns=\"" + MathMLNamespace + "\" xmlns:m=\"" + MathMLNamespace + "\">" + source + "</wrapper>")

	// create an empty MathML Element
	math = &MathML{}

	// parse the xml node
	math.root, err = xmlquery.Parse(reader)
	err = errors.Wrap(err, "xmlquery.Parse failed")
	if err != nil {
		return
	}

	// first find the MathML:semantics node
	semanticsRoot := findElementByLocalName(math.root, "semantics")
	if semanticsRoot == nil {
		return nil, errors.Errorf("No <semantics> element found in %q", math.root.OutputXML(true))
	}

	// find presentation and content
	var annotation *xmlquery.Node
	math.semantics, annotation = findPresentationAndContent(semanticsRoot)

	if math.semantics == nil {
		return nil, errors.Errorf("Did not find any non-<annotation> or non-<annotation-xml> children in %q", math.root.OutputXML(true))
	}
	if annotation == nil {
		return nil, errors.New("[ParseMathML] <semantics> element did not contain any MathML-Content <annotation-xml>")
	}

	// take its' first child
	annotation = firstChildElement(annotation)
	if annotation == nil {
		return nil, errors.New("[ParseMathML] <annotation-xml> did not contain any children")
	}

	// and update the semantics
	err = math.updateSemantics(annotation)
	err = errors.Wrap(err, "Updating semantics failed")
	return
}

func filterByAttributeLocalName(nodesIn []*xmlquery.Node, localName string, value string) (nodes []*xmlquery.Node) {
	for _, node := range nodesIn {
		for _, attr := range node.Attr {
			if attr.Name.Local == localName {
				nodes = append(nodes, node)
				break
			}
		}
	}
	return
}

// NavigateAnnotation navigates within an <annotation> element and updates the semantics accordingly
// if the xpth is invalid or an error occurs, annotation and semantics remain unchanged
func (math *MathML) NavigateAnnotation(xpth string) (err error) {
	// if we have no xpath, we have nothing to do
	// and can return immediatly
	if xpth == "" {
		return
	}

	// make sure the xpath compiles
	if _, err = xpath.Compile(xpth); err != nil {
		return
	}

	// resolve the annotation element
	annotation := xmlquery.FindOne(math.annotation, xpth)
	if annotation == nil {
		return errors.Errorf("XPath %q inside %q did not return any results", xpth, math.annotation.OutputXML(true))
	}

	// update the semantics element
	err = math.updateSemantics(annotation)
	err = errors.Wrap(err, "Updating semantics failed")

	return
}

// Copy makes a copy of this struct, allowing NavigateAnnotation() to not change the original object
func (math *MathML) Copy() *MathML {
	return &MathML{
		root:       math.root,
		semantics:  math.semantics,
		annotation: math.annotation,
	}
}

// updateSemantics updates the annotation reference and the corresponding semantics element
// if an error occurs, nothing is changed
func (math *MathML) updateSemantics(annotation *xmlquery.Node) (err error) {
	// find the xref
	xref := annotation.SelectAttr("xref")
	if xref == "" {
		return errors.Errorf("Missing xref attribute in %q", annotation.OutputXML(true))
	}

	// find the referenced element or bail out
	semantics := findElementByID(math.root, xref)
	if semantics == nil {
		return errors.Errorf("Missing <semantics> child with id %s in:\n%s", xref, math.semantics.OutputXML(true))
	}

	math.semantics = semantics
	math.annotation = annotation

	return
}

// OutputXML turns this node into an <m:math> element
// It hard-codes the "m" namespace and uses it for all elements by default.
// When no semantics element exists, it returns an <m:semantics> with only an annotation child.
func (math *MathML) OutputXML() string {
	// load semantics if they are not nil
	var semantics string
	if math.semantics != nil {
		semantics = math.semantics.OutputXML(true)
	}

	// load annotation if they are not nil
	var annotation string
	if math.annotation != nil {
		annotation = "<m:annotation-xml encoding=\"MathML-Content\">" + math.annotation.OutputXML(true) + "</m:annotation-xml>"
	}

	// build an appropriate math element
	return "<m:math xmlns=\"" + MathMLNamespace + "\" xmlns:m=\"" + MathMLNamespace + "\">" +
		"<m:semantics>" + semantics + annotation + "</m:semantics></m:math>"
}
