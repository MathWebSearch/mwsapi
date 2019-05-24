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
	semanticsRoot := xmlquery.FindOne(math.root, "//*[local-name()='semantics']")
	if semanticsRoot == nil {
		return nil, errors.New("[ParseMathML] No <semantics> found")
	}
	// then the first non-annotation node
	math.semantics = xmlquery.FindOne(semanticsRoot, "./*[not(local-name()='annotation' or local-name()='annotation-xml')]")
	if math.semantics == nil {
		return nil, errors.New("[ParseMathML] <semantics> did not contain any non-<annotation> or non-<annotation-xml> child")
	}

	// our xpath implementation does not seem to support dynamic attributes, e.g. @*[local-name()='encoding']=
	/*
		math.annotation = xmlquery.FindOne(semanticsRoot, ".//*[local-name()='annotation-xml'][*encoding='MathML-Content']/*[1]")
		if math.annotation == nil {
			return nil, errors.New("[ParseMathML] <semantics> element did not contain any MathML-Content <annotation-xml>")
		}
	*/

	// so we have a workaround
	annotations := xmlquery.Find(semanticsRoot, ".//*[local-name()='annotation-xml']")
	annotations = filterByAttributeLocalName(annotations, "encoding", "MathML-Content")
	if len(annotations) == 0 {
		return nil, errors.New("[ParseMathML] <semantics> element did not contain any MathML-Content <annotation-xml>")
	}
	math.annotation = xmlquery.FindOne(annotations[0], "./*[1]")
	if len(annotations) == 0 {
		return nil, errors.New("[ParseMathML] <annotation-xml> did not contain any children")
	}

	// and update the semantics
	err = math.updateSemantics()
	err = errors.Wrap(err, "math.updateSemantics failed")
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

// NavigateAnnotation navigates within an annotation element
// if updateSemantics is set to true, it also forcibly to semantics element using an xref attribute.
// if updateSemantics is set to false, it also tries to update the semantics element, but if updating fails, ignores the error
func (math *MathML) NavigateAnnotation(xpth string, updateSemantics bool) (err error) {
	// if we have no xpath, we have nothing to do
	// and can return immediatly
	if xpth == "" {
		return
	}

	// make sure the xpath compiles
	if _, err = xpath.Compile(xpth); err != nil {
		return
	}

	math.annotation = xmlquery.FindOne(math.annotation, xpth)
	if math.annotation == nil {
		return errors.New("[MathML.NavigateAnnotation] XPath inside <annotation-xml> did not return any results")
	}

	// update the semantics element
	err = math.updateSemantics()
	err = errors.Wrap(err, "math.updateSemantics failed")

	// if we asked to update semantics and got an error
	// ignore the error, and set the semantics to nil
	if err != nil && !updateSemantics {
		err = nil
		math.semantics = nil
	}

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

// update presentation updates the presentation element
func (math *MathML) updateSemantics() (err error) {
	// find the xref
	xref := math.annotation.SelectAttr("xref")
	if xref == "" {
		return errors.New("[MathML.updatePresentation] Missing xref attribute in <annotation-xml>")
	}

	// escape it with "s around it
	if strings.ContainsRune(xref, '"') {
		if strings.ContainsRune(xref, '\'') {
			return errors.New("[MathML.updatePresentation] xref attribute contains both single and double quote")
		}
		xref = "'" + xref + "'"
	} else {
		xref = "\"" + xref + "\""
	}

	math.semantics = xmlquery.FindOne(math.root, "//*[@xml:id="+xref+"]")
	if math.semantics == nil {
		return errors.New("[MathML.updatePresentation] Missing <semantics> child with id")
	}

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
