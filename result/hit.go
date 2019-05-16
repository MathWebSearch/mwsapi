package result

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
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

// PopulateMath populates the 'Math' property of this hit using the 'math' property of the result
func (hit *Hit) PopulateMath() error {
	// if we already have some math, return an error
	if len(hit.Math) > 0 {
		return errors.New("PopulateMath: Math already populated")
	}

	for _, mwsid := range hit.Element.MWSNumbers {
		// load the data
		data, ok := hit.Element.MWSPaths[mwsid]
		if !ok {
			return fmt.Errorf("Result %q missing path info for %d", hit.ID, mwsid)
		}

		// sort the keys in alphabetical order
		keys := make([]string, len(data))
		count := 0
		for key := range data {
			keys[count] = key
			count++
		}
		sort.Strings(keys)

		// and iterate over it
		for _, key := range keys {
			res := &MathFormula{XPath: data[key].XPath}
			res.SetURL(key)
			hit.Math = append(hit.Math, res)
		}
	}
	return nil
}
