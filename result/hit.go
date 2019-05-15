package result

import (
	"encoding/json"
	"encoding/xml"
)

// Hit represents a single Hit
type Hit struct {
	ID string `json:"id,omitempty"` // the (possibly internal) id of this hit

	URL   string `json:"url,omitempty"` // the url of the document returned
	XPath string `json:"xpath"`         // the xpath of the query term to the formulae referred to by this id

	Element *HarvestElement `json:"source,omitempty"` // the raw harvest element (if any)

	Score float64 `json:"score,omitempty"` // score of this hit

	Snippets []string `json:"snippets,omitempty"` // extracts of this hit (if any)

	Math []*MathFormula `json:"math_ids"` // math found within this hit
}

// UnmarshalJSON unmarshals a json hit
func (hit *Hit) UnmarshalJSON(bytes []byte) error {
	type marshalHit Hit // to prevent infinite recursion
	h := struct {
		XHTML string `json:"xhtml,omitempty"` // xml version of a harvest element (if any)
		*marshalHit
	}{
		marshalHit: (*marshalHit)(hit),
	}

	// unmarshal the helper
	if err := json.Unmarshal(bytes, &h); err != nil {
		return err
	}

	// no harvest element => done
	if h.XHTML == "" {
		return nil
	}

	// else unmarshal the xml as an element
	return xml.Unmarshal([]byte("<data>"+h.XHTML+"</data>"), &h.Element)
}
