package result

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/pkg/errors"
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

// innerXML is a utility struct used to extract inner xml from specific elements
type innerXML struct {
	InnerXML string `xml:",innerxml"`
}

// UnmarshalXML unmarshals xml into a harvest element
func (he *HarvestElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// try to decode the raw data or fail
	var h struct {
		XMLName xml.Name `xml:"data"`

		ID   innerXML `xml:"id"`
		Text innerXML `xml:"text"`

		Metadata innerXML `xml:"metadata"`

		Math []*MathFormula `xml:"math"`
	}
	if err := d.DecodeElement(&h, &start); err != nil {
		return errors.Wrap(err, "d.DecodeElement failed")
	}

	// start creating a harvest element
	*he = HarvestElement{
		Segment: h.ID.InnerXML,
		Text:    h.Text.InnerXML,
	}

	// load the metadata, set it to nil if undefined
	v := strings.TrimSpace(h.Metadata.InnerXML)
	if v != "" {
		// try and unmarshal it, but set it to the trimmed string on failure
		if err := json.Unmarshal([]byte(v), &he.Metadata); err != nil {
			he.Metadata = v
		}
	} else {
		he.Metadata = map[string]interface{}{}
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
