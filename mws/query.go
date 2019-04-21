package mws

import (
	"encoding/xml"

	"github.com/MathWebSearch/mwsapi/utils"
)

// Query represents a MathWebSearch Query
type Query struct {
	XMLName xml.Name `xml:"mws:query"`

	NamespaceMWS string `xml:"xmlns:mws test,attr"`

	// TODO: Define two namespaces

	From int64 `xml:"limitmin,attr"` // offset within the set of results
	Size int64 `xml:"answsize,attr"` // maximum number of results returned

	ReturnTotal  utils.BooleanYesNo `xml:"totalreq,attr"` // if true also compute the total number of elements
	OutputFormat string             `xml:"output,attr"`   // output format, "xml" or "json"

	Expressions []*Expression // the expressions that we are searching for
}

// Expression represents a single expression that is being searched for
type Expression struct {
	XMLName xml.Name `xml:"mws:expr"`  // TODO: Add some more things
	Term    string   `xml:",innerxml"` // the actual term being searched for
}
