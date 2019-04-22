package mws

import (
	"encoding/xml"

	"github.com/MathWebSearch/mwsapi/utils"
)

// RawQuery represents a MathWebSearch RawQuery
type RawQuery struct {
	From int64 `xml:"limitmin,attr"` // offset within the set of results
	Size int64 `xml:"answsize,attr"` // maximum number of results returned

	ReturnTotal  utils.BooleanYesNo `xml:"totalreq,attr"` // if true also compute the total number of elements
	OutputFormat string             `xml:"output,attr"`   // output format, "xml" or "json"

	Expressions []*Expression // the expressions that we are searching for

	XMLName xml.Name `xml:"mws:query"`
}

// Expression represents a single expression that is being searched for
type Expression struct {
	Term string `xml:",innerxml"` // the actual term being searched for

	XMLName xml.Name `xml:"mws:expr"` // TODO: Add some more things
}

// a supertype of Query that can be marshalled
type xQuery struct {
	*RawQuery

	NamespaceMWS string `xml:"xmlns:mws,attr"`
	NamespaceM   string `xml:"xmlns:m,attr"`
}

// ToXML turns a query into valid XML
func (q *RawQuery) ToXML() ([]byte, error) {
	return xml.Marshal(&xQuery{
		RawQuery:     q,
		NamespaceMWS: "http://www.mathweb.org/mws/ns",
		NamespaceM:   "http://www.w3.org/1998/Math/MathML",
	})
}
