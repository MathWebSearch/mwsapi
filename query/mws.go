package query

import (
	"encoding/xml"

	"github.com/MathWebSearch/mwsapi/utils"
)

// MWSQuery represents a user-provided MWS Query
type MWSQuery struct {
	Expressions []string // MathWebSearch Expressions to query for
	MwsIdsOnly  bool     // if set to true, use method "mws_ids", else "json"
}

// Raw turns this query  into a raw mws query
func (q *MWSQuery) Raw(from int64, size int64) *RawMWSQuery {
	// make the expressions
	exprs := make([]*MWSExpression, len(q.Expressions))
	for i, expr := range q.Expressions {
		exprs[i] = &MWSExpression{
			Term: expr,
		}
	}

	var format string
	if q.MwsIdsOnly {
		format = "mws-ids"
	} else {
		format = "json"
	}

	// and make the new raw query
	return &RawMWSQuery{
		From: from,
		Size: size,

		ReturnTotal:  true,
		OutputFormat: format,

		Expressions: exprs,
	}
}

// RawMWSQuery represents a (raw) MathWebSearch Query that is sent directly to MathWebSearch
type RawMWSQuery struct {
	From int64 `xml:"limitmin,attr"` // offset within the set of results
	Size int64 `xml:"answsize,attr"` // maximum number of results returned

	ReturnTotal  utils.BooleanYesNo `xml:"totalreq,attr"` // if true also compute the total number of elements
	OutputFormat string             `xml:"output,attr"`   // output format, "xml" or "json"

	Expressions []*MWSExpression // the expressions that we are searching for

	XMLName xml.Name `xml:"mws:query"`
}

// MWSExpression represents a single expression that is being searched for
type MWSExpression struct {
	Term    string   `xml:",innerxml"` // the actual term being searched for
	XMLName xml.Name `xml:"mws:expr"`
}

// a supertype of MWSQuery that can be marshalled
type xQuery struct {
	*RawMWSQuery

	NamespaceMWS string `xml:"xmlns:mws,attr"`
	NamespaceM   string `xml:"xmlns:m,attr"`
}

// ToXML turns a query into valid XML
func (q *RawMWSQuery) ToXML() ([]byte, error) {
	return xml.Marshal(&xQuery{
		RawMWSQuery:  q,
		NamespaceMWS: "http://www.mathweb.org/mws/ns",
		NamespaceM:   "http://www.w3.org/1998/Math/MathML",
	})
}
