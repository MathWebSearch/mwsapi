package result

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/MathWebSearch/mwsapi/utils"
)

// HarvestElement represents an <mws:data> harvest node, which can be contained both within ElasticSearch and MWS
type HarvestElement struct {
	Metadata interface{} `json:"metadata"`          // Arbitrary document metadata
	Segment  string      `json:"segment,omitempty"` // Name of the segment this document belongs to

	Text string `json:"text"` // Text of this element

	MWSPaths   map[int64]map[string]MathFormula `json:"mws_id,omitempty"`  // information about each identifier within this document
	MWSNumbers []int64                          `json:"mws_ids,omitempty"` // list of identifiers within this document

	MathSource map[string]string `json:"math"` // Source of replaced math elements within this document
}

// UnmarshalXML unmarshals xml into a harvest element
func (he *HarvestElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	// raw data
	var h struct {
		XMLName xml.Name `xml:"data"`

		ID   utils.InnerXML `xml:"id"`
		Text utils.InnerXML `xml:"text"`

		Metadata utils.InnerXML `xml:"metadata"`

		Math []*MathFormula `xml:"math"`
	}
	err = d.DecodeElement(&h, &start)
	if err != nil {
		return
	}

	// turn it into a harvestelement
	*he = HarvestElement{
		Segment: string(h.ID),
		Text:    string(h.Text),
	}

	// load the metadata, and set it to an {} if omitted
	v := strings.TrimSpace(string(h.Metadata))
	if v == "" {
		v = "{}"
	}
	err = json.Unmarshal([]byte(v), &he.Metadata)
	if err != nil {
		return
	}

	// iterate over the found math elements
	he.MathSource = make(map[string]string, len(h.Math))
	for _, math := range h.Math {
		he.MathSource[math.LocalID] = math.Source
	}

	return nil
}

// HarvestSegment represents the name of the segment
type HarvestSegment struct {
	ID string `json:"segment"`

	Hash    string `json:"hash"`    // the hash of this segment (if any)
	Touched bool   `json:"touched"` // has this segment been touched within recent changes
}
