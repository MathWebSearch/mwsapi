package utils

import (
	"encoding/xml"
	"strings"

	"github.com/antchfx/xpath"

	"github.com/antchfx/xmlquery"
)

// ResolveXPath resolves an XPath within an XHTML Fragment containing one node
func ResolveXPath(fragment string, pth string) (results []string, err error) {
	// try and compile the xpath to check if it is valid
	_, err = xpath.Compile(pth)
	if err != nil {
		return
	}

	// read the xml fragment
	reader := strings.NewReader(xml.Header + fragment)
	context, err := xmlquery.Parse(reader)
	if err != nil {
		return
	}

	// turn all the nodes into strings
	nodes := xmlquery.Find(context, pth)
	results = make([]string, len(nodes))
	for i, node := range nodes {
		results[i] = node.OutputXML(true)
	}

	return
}
